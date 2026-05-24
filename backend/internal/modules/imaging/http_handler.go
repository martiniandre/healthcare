package imaging

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/api/middleware"
	"github.com/healthcare/backend/internal/api/render"
	"github.com/healthcare/backend/internal/modules/auth"
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
	clinicalRead := middleware.RequireRoles(auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse)
	clinicalWrite := middleware.RequireRoles(auth.RoleDoctor, auth.RoleNurse)

	mux.Handle("GET /api/patients/{patientFhirId}/studies", clinicalRead(http.HandlerFunc(handler.ListPatientStudies)))
	mux.Handle("POST /api/patients/{patientFhirId}/studies", clinicalWrite(http.HandlerFunc(handler.UploadPatientStudy)))
	mux.Handle("GET /api/studies/{studyId}", clinicalRead(http.HandlerFunc(handler.GetStudy)))
}

func (handler *HTTPHandler) ListPatientStudies(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.PathValue("patientFhirId")

	studiesList, studiesErr := handler.service.ListImagingStudies(httpRequest.Context(), patientFhirID)
	if studiesErr != nil {
		slog.Error("failed to list imaging studies", "error", studiesErr, "patient_fhir_id", patientFhirID)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao carregar estudos de imagem do paciente.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, NewHTTPImagingStudyResponses(studiesList))
}

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
		slog.Error("failed to upload DICOM study", "error", createErr, "patient_fhir_id", patientFhirID)
		render.Error(httpResponseWriter, http.StatusInternalServerError, createErr.Error())
		return
	}

	render.JSON(httpResponseWriter, http.StatusCreated, NewHTTPImagingStudyResponse(createdStudy))
}

func (handler *HTTPHandler) GetStudy(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	studyIDRaw := httpRequest.PathValue("studyId")

	studyUUID, parseErr := uuid.Parse(studyIDRaw)
	if parseErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "ID de estudo inválido.")
		return
	}

	study, getStudyErr := handler.service.GetImagingStudy(httpRequest.Context(), studyUUID)
	if getStudyErr != nil {
		slog.Error("study not found", "error", getStudyErr, "study_id", studyIDRaw)
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
