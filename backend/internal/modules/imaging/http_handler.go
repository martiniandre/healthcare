package imaging

import (
	"errors"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/healthcare/backend/internal/api/middleware"
	"github.com/healthcare/backend/internal/api/render"
)

const (
	uploadReadDeadline       = 2 * time.Hour
	maximumFormFieldBytes    = 4096
	uploadFormErrorMessage   = "Erro ao processar form-data ou arquivo acima do limite permitido."
	uploadFileMissingMessage = "Arquivo DICOM ausente."
	uploadFileFieldName      = "file"
	uploadTitleFieldName     = "title"
	uploadModalityFieldName  = "modality"
)

var errFormFieldTooLarge = errors.New("form field exceeds the maximum size")

type HTTPHandler struct {
	service Service
}

func NewHTTPHandler(service Service) *HTTPHandler {
	return &HTTPHandler{
		service: service,
	}
}

func (handler *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/v1/patients/{patientFhirId}/studies", middleware.RequirePolicy("GET /api/v1/patients/{patientFhirId}/studies")(http.HandlerFunc(handler.ListPatientStudies)))
	mux.Handle("POST /api/v1/patients/{patientFhirId}/studies", middleware.RequirePolicy("POST /api/v1/patients/{patientFhirId}/studies")(http.HandlerFunc(handler.UploadPatientStudy)))
	mux.Handle("GET /api/v1/studies/{studyId}", middleware.RequirePolicy("GET /api/v1/studies/{studyId}")(http.HandlerFunc(handler.GetStudy)))
	mux.Handle("GET /api/v1/studies/{studyId}/download-url", middleware.RequirePolicy("GET /api/v1/studies/{studyId}/download-url")(http.HandlerFunc(handler.GetStudyDownloadURL)))
	mux.Handle("GET /api/v1/studies/{studyId}/instances", middleware.RequirePolicy("GET /api/v1/studies/{studyId}/instances")(http.HandlerFunc(handler.ListStudyInstances)))
}

// ListPatientStudies godoc
//
//	@Summary		List imaging studies for a patient
//	@Description	Returns all imaging studies associated with a patient
//	@Tags			imaging
//	@Accept			json
//	@Produce		json
//	@Param			patientFhirId	path	string	true	"Patient FHIR ID"
//	@Success		200				{array}	HTTPImagingStudyResponse
//	@Failure		400				{object}	map[string]string
//	@Failure		500				{object}	map[string]string
//	@Router			/patients/{patientFhirId}/studies [get]
func (handler *HTTPHandler) ListPatientStudies(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.PathValue("patientFhirId")

	studiesList, studiesErr := handler.service.ListImagingStudies(httpRequest.Context(), patientFhirID)
	if studiesErr != nil {
		slog.Error("failed to list imaging studies", "error", studiesErr, "patient_fhir_id", patientFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, studiesErr)
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, NewHTTPImagingStudyResponses(studiesList))
}

// UploadPatientStudy godoc
//
//	@Summary		Upload DICOM study or ZIP archive
//	@Description	Uploads one DICOM file, or a ZIP archive of DICOM files of the same study, for a patient. The first DICOM instance of a StudyInstanceUID creates the study; later uploads of the same StudyInstanceUID for the same patient append their instances (duplicated SOP instances are rejected). A StudyInstanceUID that already belongs to another patient is rejected with 422 DICOM_PATIENT_MISMATCH. In a ZIP, a patient mismatch in any instance rejects the whole request and nothing is stored; other per-instance problems are listed in instancesRejected while the remaining instances are stored, and an archive with no accepted instance returns 422 DICOM_NO_VALID_INSTANCES. Archives must hold one study, reject unsafe entry paths and links (422 ZIP_UNSAFE_ENTRY), and exceed size, entry count or compression limits with 413 ZIP_LIMITS_EXCEEDED. Content sent as a ZIP that is not a ZIP is 422 ZIP_INVALID. A single DICOM file that was already uploaded is rejected with 409 DICOM_DUPLICATE_INSTANCE. Title and modality are optional fallbacks; the DICOM modality wins when both are present. The multipart form fields title and modality must be sent BEFORE the file part; fields that arrive after the file part are ignored. The read deadline of this route is extended to 2 hours because the global 30 second read timeout would cut large uploads; the request body is streamed to storage and never spooled to disk as a whole form, except ZIP archives which are spooled to a temporary file with size limits.
//	@Tags			imaging
//	@Accept			mpfd
//	@Produce		json
//	@Param			patientFhirId	path		string	true	"Patient FHIR ID"
//	@Param			title			formData	string	false	"Study title (derived from the DICOM header when empty). Must precede the file part."
//	@Param			modality		formData	string	false	"Fallback modality (e.g., CT, MR, XR) used only when the DICOM header has none. Must precede the file part."
//	@Param			file			formData	file	true	"DICOM file or ZIP archive of DICOM files (.zip)"
//	@Success		201				{object}	HTTPUploadStudyResponse
//	@Failure		400				{object}	map[string]string
//	@Failure		404				{object}	map[string]string
//	@Failure		409				{object}	map[string]string
//	@Failure		413				{object}	map[string]string
//	@Failure		422				{object}	map[string]string
//	@Failure		429				{object}	map[string]string
//	@Failure		500				{object}	map[string]string
//	@Router			/patients/{patientFhirId}/studies [post]
func (handler *HTTPHandler) UploadPatientStudy(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.PathValue("patientFhirId")
	_ = http.NewResponseController(httpResponseWriter).SetReadDeadline(time.Now().Add(uploadReadDeadline))

	multipartReader, readerError := httpRequest.MultipartReader()
	if readerError != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, uploadFormErrorMessage)
		return
	}

	formFields := make(map[string]string)
	for {
		part, partError := multipartReader.NextPart()
		if errors.Is(partError, io.EOF) {
			render.Error(httpResponseWriter, http.StatusBadRequest, uploadFileMissingMessage)
			return
		}
		if partError != nil {
			render.Error(httpResponseWriter, http.StatusBadRequest, uploadFormErrorMessage)
			return
		}

		if part.FormName() == uploadFileFieldName && part.FileName() != "" {
			handler.streamUploadedFile(httpResponseWriter, httpRequest, patientFhirID, formFields, part)
			return
		}

		if fieldName := part.FormName(); fieldName == uploadTitleFieldName || fieldName == uploadModalityFieldName {
			fieldValue, valueError := readFormFieldValue(part)
			if valueError != nil {
				render.Error(httpResponseWriter, http.StatusBadRequest, uploadFormErrorMessage)
				return
			}
			formFields[fieldName] = fieldValue
		}
	}
}

func (handler *HTTPHandler) streamUploadedFile(httpResponseWriter http.ResponseWriter, httpRequest *http.Request, patientFhirID string, formFields map[string]string, filePart *multipart.Part) {
	uploadResult, createErr := handler.service.UploadDICOMStream(httpRequest.Context(), UploadDICOMInput{
		PatientFhirID: patientFhirID,
		Title:         formFields[uploadTitleFieldName],
		Modality:      formFields[uploadModalityFieldName],
		FileName:      filePart.FileName(),
	}, filePart)
	if createErr != nil {
		slog.Error("failed to upload DICOM study", "error", createErr, "patient_fhir_id", patientFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, createErr)
		return
	}

	render.JSON(httpResponseWriter, http.StatusCreated, NewHTTPUploadStudyResponse(uploadResult))
}

func readFormFieldValue(part *multipart.Part) (string, error) {
	fieldBytes, readError := io.ReadAll(io.LimitReader(part, maximumFormFieldBytes+1))
	if readError != nil {
		return "", readError
	}
	if len(fieldBytes) > maximumFormFieldBytes {
		return "", errFormFieldTooLarge
	}
	return string(fieldBytes), nil
}

// GetStudy godoc
//
//	@Summary		Get imaging study details
//	@Description	Returns details of a specific imaging study including download URL. The download URL is omitted for studies with more than one instance; use the instances endpoint for those.
//	@Tags			imaging
//	@Accept			json
//	@Produce		json
//	@Param			studyId	path	string	true	"Study UUID"
//	@Success		200		{object}	HTTPImagingStudyResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Router			/studies/{studyId} [get]
func (handler *HTTPHandler) GetStudy(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	studyID := httpRequest.PathValue("studyId")

	study, getStudyErr := handler.service.GetImagingStudy(httpRequest.Context(), studyID)
	if getStudyErr != nil {
		slog.Error("failed to get imaging study", "error", getStudyErr, "study_id", studyID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, getStudyErr)
		return
	}

	downloadURL, _, downloadErr := handler.service.GetDownloadURL(httpRequest.Context(), studyID)
	if downloadErr != nil {
		slog.Warn("failed to generate download URL", "study_id", studyID, "error", downloadErr)
		downloadURL = ""
	}

	render.JSON(httpResponseWriter, http.StatusOK, NewHTTPImagingStudyResponseWithDownloadURL(study, downloadURL))
}

// GetStudyDownloadURL godoc
//
//	@Summary		Get study download URL
//	@Description	Returns a presigned download URL (valid 15 minutes) for a single-instance study. Studies with more than one instance return 409 MULTIPLE_INSTANCES; use GET /studies/{studyId}/instances instead.
//	@Tags			imaging
//	@Accept			json
//	@Produce		json
//	@Param			studyId	path		string	true	"Study UUID"
//	@Success		200		{object}	HTTPStudyDownloadURLResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		409		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/studies/{studyId}/download-url [get]
func (handler *HTTPHandler) GetStudyDownloadURL(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	studyID := httpRequest.PathValue("studyId")

	downloadURL, expiresAt, downloadErr := handler.service.GetDownloadURL(httpRequest.Context(), studyID)
	if downloadErr != nil {
		slog.Error("failed to issue study download URL", "error", downloadErr, "study_id", studyID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, downloadErr)
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, HTTPStudyDownloadURLResponse{
		DownloadURL: downloadURL,
		ExpiresAt:   expiresAt.Format(time.RFC3339),
	})
}

// ListStudyInstances godoc
//
//	@Summary		List study instances with download URLs
//	@Description	Returns every stored DICOM instance of a study with a presigned download URL (valid 15 minutes) for each one. Each download is audited.
//	@Tags			imaging
//	@Accept			json
//	@Produce		json
//	@Param			studyId	path		string	true	"Study UUID"
//	@Success		200		{array}		HTTPStudyInstanceResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/studies/{studyId}/instances [get]
func (handler *HTTPHandler) ListStudyInstances(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	studyID := httpRequest.PathValue("studyId")

	downloads, listErr := handler.service.ListStudyInstanceDownloads(httpRequest.Context(), studyID)
	if listErr != nil {
		slog.Error("failed to list study instances", "error", listErr, "study_id", studyID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, listErr)
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, NewHTTPStudyInstanceResponses(downloads))
}
