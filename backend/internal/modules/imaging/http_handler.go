package imaging

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/api/middleware"
	"github.com/healthcare/backend/internal/api/render"
	"github.com/healthcare/backend/internal/shared/role"
)

type HTTPHandler struct {
	service Service
}

func NewHTTPHandler(service Service) *HTTPHandler {
	return &HTTPHandler{
		service: service,
	}
}

func (handler *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	clinicalRead := middleware.RequireRoles(role.RoleAdmin, role.RoleDoctor, role.RoleNurse)
	clinicalWrite := middleware.RequireRoles(role.RoleAdmin, role.RoleDoctor, role.RoleNurse)

	mux.Handle("GET /api/v1/patients/{patientFhirId}/studies", clinicalRead(http.HandlerFunc(handler.ListPatientStudies)))
	mux.Handle("POST /api/v1/patients/{patientFhirId}/studies", clinicalWrite(http.HandlerFunc(handler.UploadPatientStudy)))
	mux.Handle("GET /api/v1/studies/{studyId}", clinicalRead(http.HandlerFunc(handler.GetStudy)))
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
//	@Failure		500				{object}	map[string]string
//	@Router			/patients/{patientFhirId}/studies [get]
func (handler *HTTPHandler) ListPatientStudies(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.PathValue("patientFhirId")

	studiesList, studiesErr := handler.service.ListImagingStudies(httpRequest.Context(), patientFhirID)
	if studiesErr != nil {
		slog.Error("failed to list imaging studies", "error", studiesErr, "patient_fhir_id", patientFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao carregar estudos de imagem do paciente.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, NewHTTPImagingStudyResponses(studiesList))
}

// UploadPatientStudy godoc
//
//	@Summary		Upload DICOM study
//	@Description	Uploads a DICOM imaging study for a patient with metadata via multipart form
//	@Tags			imaging
//	@Accept			mpfd
//	@Produce		json
//	@Param			patientFhirId	path	string	true	"Patient FHIR ID"
//	@Param			title			formData	string	true	"Study title"
//	@Param			modality		formData	string	true	"Modality (e.g., CT, MR, XR)"
//	@Param			file			formData	file	true	"DICOM file"
//	@Success		201				{object}	HTTPImagingStudyResponse
//	@Failure		400				{object}	map[string]string
//	@Failure		500				{object}	map[string]string
//	@Router			/patients/{patientFhirId}/studies [post]
func (handler *HTTPHandler) UploadPatientStudy(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.PathValue("patientFhirId")

	httpRequest.Body = http.MaxBytesReader(httpResponseWriter, httpRequest.Body, MaxDICOMUploadBytes)
	if parseErr := httpRequest.ParseMultipartForm(10 << 20); parseErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Erro ao processar form-data ou arquivo acima do limite permitido.")
		return
	}
	defer httpRequest.MultipartForm.RemoveAll()

	title := httpRequest.FormValue("title")
	modality := httpRequest.FormValue("modality")
	if title == "" || modality == "" {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Campos title e modality são obrigatórios.")
		return
	}

	file, _, fileErr := httpRequest.FormFile("file")
	if fileErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Arquivo DICOM ausente.")
		return
	}
	defer file.Close()

	createdStudy, createErr := handler.service.UploadDICOMStream(httpRequest.Context(), patientFhirID, title, modality, file)
	if createErr != nil {
		slog.Error("failed to upload DICOM study", "error", createErr, "patient_fhir_id", patientFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, createErr.Error())
		return
	}

	render.JSON(httpResponseWriter, http.StatusCreated, NewHTTPImagingStudyResponse(createdStudy))
}

// GetStudy godoc
//
//	@Summary		Get imaging study details
//	@Description	Returns details of a specific imaging study including download URL
//	@Tags			imaging
//	@Accept			json
//	@Produce		json
//	@Param			studyId	path	string	true	"Study UUID"
//	@Success		200		{object}	HTTPImagingStudyResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Router			/studies/{studyId} [get]
func (handler *HTTPHandler) GetStudy(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	studyIDRaw := httpRequest.PathValue("studyId")

	studyUUID, parseErr := uuid.Parse(studyIDRaw)
	if parseErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "ID de estudo inválido.")
		return
	}

	study, getStudyErr := handler.service.GetImagingStudy(httpRequest.Context(), studyUUID)
	if getStudyErr != nil {
		slog.Error("study not found", "error", getStudyErr, "study_id", studyIDRaw, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusNotFound, "Estudo não encontrado.")
		return
	}

	downloadURL, _, downloadErr := handler.service.GetDownloadURL(httpRequest.Context(), studyUUID)
	if downloadErr != nil {
		slog.Warn("failed to generate download URL", "study_id", studyUUID, "error", downloadErr)
		downloadURL = ""
	}

	render.JSON(httpResponseWriter, http.StatusOK, NewHTTPImagingStudyResponseWithDownloadURL(study, downloadURL))
}
