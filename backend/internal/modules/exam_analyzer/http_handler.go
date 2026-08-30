package exam_analyzer

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/api/middleware"
	"github.com/healthcare/backend/internal/api/render"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/ctxkeys"
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
	mux.Handle("GET /api/v1/exam-analyses", middleware.RequirePolicy("GET /api/v1/exam-analyses")(http.HandlerFunc(handler.ListAnalyses)))
	mux.Handle("POST /api/v1/exam-analyses", middleware.RequirePolicy("POST /api/v1/exam-analyses")(http.HandlerFunc(handler.CreateAnalysis)))
	mux.Handle("GET /api/v1/exam-analyses/{analysisId}", middleware.RequirePolicy("GET /api/v1/exam-analyses/{analysisId}")(http.HandlerFunc(handler.GetAnalysis)))
	mux.Handle("DELETE /api/v1/exam-analyses/{analysisId}", middleware.RequirePolicy("DELETE /api/v1/exam-analyses/{analysisId}")(http.HandlerFunc(handler.DeleteAnalysis)))
}

// ListAnalyses godoc
//
//	@Summary		List exam analyses
//	@Description	Returns all exam analyses, optionally filtered by patient FHIR ID
//	@Tags			exam_analyzer
//	@Accept			json
//	@Produce		json
//	@Param			patientFhirId	query	string	false	"Filter by patient FHIR ID"
//	@Success		200				{array}	ExamAnalysisResponse
//	@Failure		500				{object}	map[string]string
//	@Router			/exam-analyses [get]
func (handler *HTTPHandler) ListAnalyses(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.URL.Query().Get("patientFhirId")
	var filterPatient *string
	if patientFhirID != "" {
		filterPatient = &patientFhirID
	}

	analysesList, listError := handler.repository.ListAnalyses(httpRequest.Context(), filterPatient)
	if listError != nil {
		slog.Error("failed to list exam analyses", "error", listError, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao listar análises de exames.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, analysesList)
}

// CreateAnalysis godoc
//
//	@Summary		Create exam analysis
//	@Description	Uploads a medical exam file for AI analysis with consent tracking and optional anonymization
//	@Tags			exam_analyzer
//	@Accept			mpfd
//	@Produce		json
//	@Param			patientFhirId	formData	string	false	"Patient FHIR ID"
//	@Param			consent			formData	string	true	"Patient consent (must be 'true')"
//	@Param			anonymize		formData	string	false	"Anonymize file (true/false)"
//	@Param			acknowledgeExternalProcessing	formData	string	false	"Required true when anonymize is false"
//	@Param			file			formData	file	true	"Medical exam file"
//	@Success		201				{object}	ExamAnalysisResponse
//	@Failure		400				{object}	map[string]string
//	@Failure		422				{object}	map[string]string
//	@Failure		500				{object}	map[string]string
//	@Router			/exam-analyses [post]
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

	if !isAnonymized && httpRequest.FormValue("acknowledgeExternalProcessing") != "true" {
		render.ErrorFromAppError(httpResponseWriter, apperrors.ErrDeidentificationRequired)
		return
	}

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

	analysisID := uuid.New()
	if !isSupportedExamExtension(fileHeader.Filename) {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Formato de arquivo não suportado. Envie PNG, JPG ou PDF.")
		return
	}

	sanitizedFileName := unsafeExamPathChars.ReplaceAllString(strings.TrimSpace(filepath.Base(fileHeader.Filename)), "_")
	if sanitizedFileName == "" {
		sanitizedFileName = analysisID.String() + examFileExtension(fileHeader.Filename)
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
		FileName:         sanitizedFileName,
		Status:           "pending",
		AnalysisResponse: json.RawMessage(defaultResponse),
		ConsentGiven:     true,
		Anonymized:       isAnonymized,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if saveErr := handler.service.CreateAnalysis(httpRequest.Context(), newAnalysisRecord, file); saveErr != nil {
		slog.Error("failed to save analysis metadata", "error", saveErr, "analysis_id", analysisID, "request_id", middleware.GetRequestID(httpRequest.Context()))
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

// GetAnalysis godoc
//
//	@Summary		Get exam analysis
//	@Description	Returns details of a specific exam analysis by ID
//	@Tags			exam_analyzer
//	@Accept			json
//	@Produce		json
//	@Param			analysisId	path	string	true	"Analysis UUID"
//	@Success		200			{object}	ExamAnalysisResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Router			/exam-analyses/{analysisId} [get]
func (handler *HTTPHandler) GetAnalysis(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	analysisIDRaw := httpRequest.PathValue("analysisId")

	analysisID, parseErr := uuid.Parse(analysisIDRaw)
	if parseErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "ID de análise inválido.")
		return
	}

	analysisRecord, fetchErr := handler.repository.GetAnalysis(httpRequest.Context(), analysisID)
	if fetchErr != nil {
		slog.Error("analysis not found", "error", fetchErr, "analysis_id", analysisIDRaw, "request_id", middleware.GetRequestID(httpRequest.Context()))
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

// DeleteAnalysis godoc
//
//	@Summary		Delete exam analysis
//	@Description	Permanently deletes an exam analysis and its associated file
//	@Tags			exam_analyzer
//	@Accept			json
//	@Produce		json
//	@Param			analysisId	path	string	true	"Analysis UUID"
//	@Success		200			{object}	DeleteAnalysisResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/exam-analyses/{analysisId} [delete]
func (handler *HTTPHandler) DeleteAnalysis(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	analysisIDRaw := httpRequest.PathValue("analysisId")

	analysisID, parseErr := uuid.Parse(analysisIDRaw)
	if parseErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "ID de análise inválido.")
		return
	}

	if deleteErr := handler.service.DeleteAnalysis(httpRequest.Context(), analysisID); deleteErr != nil {
		slog.Error("failed to delete analysis", "error", deleteErr, "analysis_id", analysisIDRaw, "request_id", middleware.GetRequestID(httpRequest.Context()))
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

type ExamAnalysisResponse struct {
	ID               string `json:"id"`
	UserID           string `json:"user_id"`
	PatientFhirID    string `json:"patient_fhir_id"`
	ExamType         string `json:"exam_type"`
	FileName         string `json:"file_name"`
	Status           string `json:"status"`
	AnalysisResponse string `json:"analysis_response"`
	ConsentGiven     bool   `json:"consent_given"`
	Anonymized       bool   `json:"anonymized"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

type DeleteAnalysisResponse struct {
	Success string `json:"success"`
}
