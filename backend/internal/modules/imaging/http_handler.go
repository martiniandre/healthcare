package imaging

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/auth"
)

type HTTPAuthValidator func(http.ResponseWriter, *http.Request, []auth.Role) (context.Context, bool)

type HTTPHandler struct {
	service      Service
	validateAuth HTTPAuthValidator
}

func NewHTTPHandler(service Service, validateAuth HTTPAuthValidator) *HTTPHandler {
	return &HTTPHandler{
		service:      service,
		validateAuth: validateAuth,
	}
}

func (handler *HTTPHandler) HandlePatientStudies(httpResponseWriter http.ResponseWriter, httpRequest *http.Request, patientFHIRID string) {
	if httpRequest.Method == http.MethodGet {
		contextWithValues, authIsOk := handler.validateAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse})
		if !authIsOk {
			return
		}

		studiesList, studiesErr := handler.service.ListImagingStudies(contextWithValues, patientFHIRID)
		if studiesErr != nil {
			httpResponseWriter.Header().Set("Content-Type", "application/json")
			httpResponseWriter.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao carregar estudos de imagem do paciente."})
			return
		}

		httpResponseWriter.Header().Set("Content-Type", "application/json")
		httpResponseWriter.WriteHeader(http.StatusOK)
		json.NewEncoder(httpResponseWriter).Encode(NewHTTPImagingStudyResponses(studiesList))
		return
	}

	if httpRequest.Method == http.MethodPost {
		contextWithValues, authIsOk := handler.validateAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleDoctor, auth.RoleNurse})
		if !authIsOk {
			return
		}

		httpRequest.Body = http.MaxBytesReader(httpResponseWriter, httpRequest.Body, MaxDICOMUploadBytes)
		if parseErr := httpRequest.ParseMultipartForm(10 << 20); parseErr != nil {
			httpResponseWriter.Header().Set("Content-Type", "application/json")
			httpResponseWriter.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao processar form-data ou arquivo acima do limite permitido."})
			return
		}
		defer httpRequest.MultipartForm.RemoveAll()

		title := httpRequest.FormValue("title")
		modality := httpRequest.FormValue("modality")
		if title == "" || modality == "" {
			httpResponseWriter.Header().Set("Content-Type", "application/json")
			httpResponseWriter.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Campos title e modality são obrigatórios."})
			return
		}

		file, _, fileErr := httpRequest.FormFile("file")
		if fileErr != nil {
			httpResponseWriter.Header().Set("Content-Type", "application/json")
			httpResponseWriter.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Arquivo DICOM ausente."})
			return
		}
		defer file.Close()

		createdStudy, createErr := handler.service.UploadDICOMStream(contextWithValues, patientFHIRID, title, modality, file)
		if createErr != nil {
			httpResponseWriter.Header().Set("Content-Type", "application/json")
			httpResponseWriter.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": createErr.Error()})
			return
		}

		httpResponseWriter.Header().Set("Content-Type", "application/json")
		httpResponseWriter.WriteHeader(http.StatusCreated)
		json.NewEncoder(httpResponseWriter).Encode(NewHTTPImagingStudyResponse(createdStudy))
		return
	}

	http.Error(httpResponseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func (handler *HTTPHandler) HandleStudy(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	remainingPath := strings.TrimPrefix(httpRequest.URL.Path, "/api/studies/")
	if remainingPath == "" {
		httpResponseWriter.Header().Set("Content-Type", "application/json")
		httpResponseWriter.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "ID do estudo ausente."})
		return
	}

	if httpRequest.Method != http.MethodGet {
		http.Error(httpResponseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	contextWithValues, authIsOk := handler.validateAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse})
	if !authIsOk {
		return
	}

	studyUUID, parseErr := uuid.Parse(remainingPath)
	if parseErr != nil {
		httpResponseWriter.Header().Set("Content-Type", "application/json")
		httpResponseWriter.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "ID de estudo inválido."})
		return
	}

	study, getStudyErr := handler.service.GetImagingStudy(contextWithValues, studyUUID)
	if getStudyErr != nil {
		httpResponseWriter.Header().Set("Content-Type", "application/json")
		httpResponseWriter.WriteHeader(http.StatusNotFound)
		json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Estudo não encontrado."})
		return
	}

	downloadURL, _, downloadErr := handler.service.GetDownloadURL(contextWithValues, studyUUID)
	if downloadErr != nil {
		slog.Warn("Failed to generate download URL", "studyID", studyUUID, "error", downloadErr)
		downloadURL = ""
	}

	httpResponseWriter.Header().Set("Content-Type", "application/json")
	httpResponseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(httpResponseWriter).Encode(NewHTTPImagingStudyResponseWithDownloadURL(study, downloadURL))
}
