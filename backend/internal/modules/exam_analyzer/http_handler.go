package exam_analyzer

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/api/middleware"
	"github.com/healthcare/backend/internal/api/render"
	"github.com/healthcare/backend/internal/shared/ctxkeys"
	"github.com/healthcare/backend/internal/shared/role"
)

type HTTPHandler struct {
	repository Repository
	service    Service
	worker     *Worker
}

func NewHTTPHandler(repository Repository, service Service, worker *Worker) *HTTPHandler {
	return &HTTPHandler{
		repository: repository,
		service:    service,
		worker:     worker,
	}
}

func (handler *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	clinicalRead := middleware.RequireRoles(role.RoleAdmin, role.RoleDoctor, role.RoleNurse)
	clinicalDelete := middleware.RequireRoles(role.RoleAdmin, role.RoleDoctor)

	mux.Handle("GET /api/exam-analyses", clinicalRead(http.HandlerFunc(handler.ListAnalyses)))
	mux.Handle("POST /api/exam-analyses", clinicalRead(http.HandlerFunc(handler.CreateAnalysis)))
	mux.Handle("GET /api/exam-analyses/{analysisId}", clinicalRead(http.HandlerFunc(handler.GetAnalysis)))
	mux.Handle("DELETE /api/exam-analyses/{analysisId}", clinicalDelete(http.HandlerFunc(handler.DeleteAnalysis)))
}

func (handler *HTTPHandler) ListAnalyses(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.URL.Query().Get("patientFhirId")
	var filterPatient *string
	if patientFhirID != "" {
		filterPatient = &patientFhirID
	}

	analysesList, listError := handler.repository.ListAnalyses(httpRequest.Context(), filterPatient)
	if listError != nil {
		slog.Error("failed to list exam analyses", "error", listError)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao listar análises de exames.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, analysesList)
}

func (handler *HTTPHandler) CreateAnalysis(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	httpRequest.Body = http.MaxBytesReader(httpResponseWriter, httpRequest.Body, 15<<20)
	if parseErr := httpRequest.ParseMultipartForm(15 << 20); parseErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Arquivo enviado excede o limite permitido ou form inválido.")
		return
	}
	defer httpRequest.MultipartForm.RemoveAll()

	consentValue := httpRequest.FormValue("consent")
	if consentValue != "true" {
		render.Error(httpResponseWriter, http.StatusBadRequest, "O consentimento do paciente é obrigatório para processamento clínico.")
		return
	}

	anonymizeValue := httpRequest.FormValue("anonymize")
	isAnonymized := anonymizeValue == "true"

	patientFhirID := httpRequest.FormValue("patientFhirId")
	var targetPatient *string
	if patientFhirID != "" {
		targetPatient = &patientFhirID
	}

	file, fileHeader, fileErr := httpRequest.FormFile("file")
	if fileErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Arquivo de exame médico ausente.")
		return
	}
	defer file.Close()

	uploadDirectory := filepath.Join("tmp", "exam_uploads")
	if makeDirErr := os.MkdirAll(uploadDirectory, 0755); makeDirErr != nil {
		slog.Error("failed to create upload directory", "error", makeDirErr, "path", uploadDirectory)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao inicializar pasta temporária de uploads.")
		return
	}

	analysisID := uuid.New()
	fileExtension := filepath.Ext(fileHeader.Filename)

	savedFileName := fileHeader.Filename
	if isAnonymized {
		savedFileName = "anonymized_" + analysisID.String() + fileExtension
	}

	destinationPath := filepath.Join(uploadDirectory, analysisID.String()+fileExtension)
	destinationFile, createErr := os.Create(destinationPath)
	if createErr != nil {
		slog.Error("failed to create destination file", "error", createErr, "path", destinationPath)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Falha ao criar arquivo no destino temporário.")
		return
	}
	defer destinationFile.Close()

	if _, copyErr := io.Copy(destinationFile, file); copyErr != nil {
		slog.Error("failed to write file to disk", "error", copyErr, "path", destinationPath)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Falha ao gravar arquivo em disco.")
		return
	}

	userIDStr, _ := httpRequest.Context().Value(ctxkeys.UserIDKey).(string)
	userRoleStr, _ := httpRequest.Context().Value(ctxkeys.RoleKey).(string)

	var parsedUserID *uuid.UUID
	if userIDStr != "" {
		parsedID, parseUUIDErr := uuid.Parse(userIDStr)
		if parseUUIDErr == nil {
			parsedUserID = &parsedID
		}
	}

	defaultResponse, _ := json.Marshal(map[string]string{
		"status": "pending",
	})

	newAnalysisRecord := &ExamAnalysis{
		ID:               analysisID,
		UserID:           parsedUserID,
		PatientFhirID:    targetPatient,
		ExamType:         nil,
		FileName:         savedFileName,
		FilePath:         destinationPath,
		Status:           "pending",
		AnalysisResponse: json.RawMessage(defaultResponse),
		ConsentGiven:     true,
		Anonymized:       isAnonymized,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if saveErr := handler.repository.CreateAnalysis(httpRequest.Context(), newAnalysisRecord); saveErr != nil {
		slog.Error("failed to save analysis metadata", "error", saveErr, "analysis_id", analysisID)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Falha ao salvar metadados da análise.")
		return
	}

	operatorInfo := userRoleStr
	if userIDStr != "" {
		operatorInfo += " (" + userIDStr + ")"
	}

	auditDetail := "File successfully uploaded and queued for processing"
	if isAnonymized {
		auditDetail += " (Anonymization enabled)"
	}

	newAuditLog := &ExamAnalysisAuditLog{
		ID:          uuid.New(),
		AnalysisID:  &analysisID,
		ActionType:  "upload",
		PerformedBy: operatorInfo,
		IPAddress:   nil,
		Details:     &auditDetail,
		CreatedAt:   time.Now(),
	}
	_ = handler.repository.CreateAuditLog(httpRequest.Context(), newAuditLog)

	handler.worker.SubmitJob(analysisID)

	render.JSON(httpResponseWriter, http.StatusCreated, newAnalysisRecord)
}

func (handler *HTTPHandler) GetAnalysis(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	analysisIDRaw := httpRequest.PathValue("analysisId")

	analysisID, parseErr := uuid.Parse(analysisIDRaw)
	if parseErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "ID de análise inválido.")
		return
	}

	analysisRecord, fetchErr := handler.repository.GetAnalysis(httpRequest.Context(), analysisID)
	if fetchErr != nil {
		slog.Error("analysis not found", "error", fetchErr, "analysis_id", analysisIDRaw)
		render.Error(httpResponseWriter, http.StatusNotFound, "Análise de exame não encontrada.")
		return
	}

	userIDStr, _ := httpRequest.Context().Value(ctxkeys.UserIDKey).(string)
	userRoleStr, _ := httpRequest.Context().Value(ctxkeys.RoleKey).(string)
	operatorInfo := userRoleStr
	if userIDStr != "" {
		operatorInfo += " (" + userIDStr + ")"
	}

	auditMessage := "Exam analysis details accessed by medical staff"
	auditRecord := &ExamAnalysisAuditLog{
		ID:          uuid.New(),
		AnalysisID:  &analysisID,
		ActionType:  "view",
		PerformedBy: operatorInfo,
		IPAddress:   nil,
		Details:     &auditMessage,
		CreatedAt:   time.Now(),
	}
	_ = handler.repository.CreateAuditLog(httpRequest.Context(), auditRecord)

	render.JSON(httpResponseWriter, http.StatusOK, analysisRecord)
}

func (handler *HTTPHandler) DeleteAnalysis(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	analysisIDRaw := httpRequest.PathValue("analysisId")

	analysisID, parseErr := uuid.Parse(analysisIDRaw)
	if parseErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "ID de análise inválido.")
		return
	}

	analysisRecord, fetchErr := handler.repository.GetAnalysis(httpRequest.Context(), analysisID)
	if fetchErr != nil {
		slog.Error("analysis not found for deletion", "error", fetchErr, "analysis_id", analysisIDRaw)
		render.Error(httpResponseWriter, http.StatusNotFound, "Análise de exame não encontrada para exclusão.")
		return
	}

	if analysisRecord.FilePath != "deleted" {
		if _, statErr := os.Stat(analysisRecord.FilePath); statErr == nil {
			_ = os.Remove(analysisRecord.FilePath)
		}
	}

	if deleteErr := handler.repository.DeleteAnalysis(httpRequest.Context(), analysisID); deleteErr != nil {
		slog.Error("failed to delete analysis", "error", deleteErr, "analysis_id", analysisIDRaw)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Falha ao remover análise de exame do banco de dados.")
		return
	}

	userIDStr, _ := httpRequest.Context().Value(ctxkeys.UserIDKey).(string)
	userRoleStr, _ := httpRequest.Context().Value(ctxkeys.RoleKey).(string)
	operatorInfo := userRoleStr
	if userIDStr != "" {
		operatorInfo += " (" + userIDStr + ")"
	}

	auditMessage := "Exam analysis physically deleted and purged by user action"
	auditRecord := &ExamAnalysisAuditLog{
		ID:          uuid.New(),
		AnalysisID:  nil,
		ActionType:  "delete",
		PerformedBy: operatorInfo,
		IPAddress:   nil,
		Details:     &auditMessage,
		CreatedAt:   time.Now(),
	}
	_ = handler.repository.CreateAuditLog(httpRequest.Context(), auditRecord)

	render.JSON(httpResponseWriter, http.StatusOK, map[string]string{"success": "Análise e arquivo excluídos com sucesso."})
}
