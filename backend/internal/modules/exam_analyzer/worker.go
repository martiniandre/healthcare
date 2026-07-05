package exam_analyzer

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/eventbus"
)

type Worker struct {
	repository Repository
	service    Service
	jobChannel chan uuid.UUID
	stopSignal chan struct{}
	eventBus   eventbus.Bus
}

func NewWorker(repository Repository, service Service, eventBus eventbus.Bus) *Worker {
	return &Worker{
		repository: repository,
		service:    service,
		eventBus:   eventBus,
		jobChannel: make(chan uuid.UUID, 100),
		stopSignal: make(chan struct{}),
	}
}

func (worker *Worker) Start(ctx context.Context) {
	slog.Info("Exam Analyzer Background Worker initialized")

	cleanupTicker := time.NewTicker(5 * time.Minute)
	defer cleanupTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-worker.stopSignal:
			return
		case jobID := <-worker.jobChannel:
			worker.processAnalysisJob(ctx, jobID)
		case <-cleanupTicker.C:
			worker.executeAutoCleanup(ctx)
		}
	}
}

func (worker *Worker) SubmitJob(analysisID uuid.UUID) {
	worker.jobChannel <- analysisID
}

func (worker *Worker) Stop() {
	close(worker.stopSignal)
}

func (worker *Worker) processAnalysisJob(ctx context.Context, analysisID uuid.UUID) {
	analysisRecord, err := worker.repository.GetAnalysis(ctx, analysisID)
	if err != nil {
		slog.Error("Failed to fetch analysis record for processing", "analysisID", analysisID, "error", err)
		return
	}

	analysisRecord.Status = "processing"
	analysisRecord.UpdatedAt = time.Now()
	if updateErr := worker.repository.UpdateAnalysis(ctx, analysisRecord); updateErr != nil {
		slog.Error("Failed to update status to processing", "analysisID", analysisID, "error", updateErr)
		return
	}

	var analysisResponse *MedicalAnalysisResponse
	var statusResult string
	maxRetries := 3

	for attempt := 1; attempt <= maxRetries; attempt++ {
		analysisResponse, statusResult, err = worker.service.AnalyzeExamFile(ctx, analysisRecord.FilePath, analysisRecord.FileName)
		if err == nil {
			break
		}

		slog.Warn("Error during medical exam analysis execution, retrying...", "analysisID", analysisID, "attempt", attempt, "error", err)
		if attempt < maxRetries {
			time.Sleep(time.Duration(attempt*2) * time.Second)
		}
	}

	if err != nil {
		slog.Error("Max retries reached. Medical exam analysis failed", "analysisID", analysisID, "error", err)
		analysisRecord.Status = "failed"
		analysisRecord.UpdatedAt = time.Now()
		_ = worker.repository.UpdateAnalysis(ctx, analysisRecord)

		auditMessage := "Analysis moved to DLQ (failed state) after max retries"
		auditRecord := &ExamAnalysisAuditLog{
			ID:          uuid.New(),
			AnalysisID:  &analysisID,
			ActionType:  "failed",
			PerformedBy: "SYSTEM_WORKER",
			Details:     &auditMessage,
			CreatedAt:   time.Now(),
		}
		_ = worker.repository.CreateAuditLog(ctx, auditRecord)
		return
	}

	if statusResult == "insufficient_data" {
		analysisRecord.Status = "insufficient_data"
		insufficientResponseBytes, _ := json.Marshal(map[string]string{
			"status":  "insufficient_data",
			"message": "Não foi possível gerar análise confiável devido à qualidade ou ilegibilidade do arquivo enviado.",
		})
		analysisRecord.AnalysisResponse = json.RawMessage(insufficientResponseBytes)
	} else {
		analysisRecord.Status = "completed"
		examTypeString := analysisResponse.ExamType
		analysisRecord.ExamType = &examTypeString
		responseBytes, marshalErr := json.Marshal(analysisResponse)
		if marshalErr != nil {
			slog.Error("Failed to marshal analysis response", "analysisID", analysisID, "error", marshalErr)
			analysisRecord.Status = "failed"
		} else {
			analysisRecord.AnalysisResponse = json.RawMessage(responseBytes)
		}
	}

	analysisRecord.UpdatedAt = time.Now()
	if updateErr := worker.repository.UpdateAnalysis(ctx, analysisRecord); updateErr != nil {
		slog.Error("Failed to save final analysis results to database", "analysisID", analysisID, "error", updateErr)
		return
	}

	auditDetail := "Automatic analysis execution completed successfully"
	if statusResult == "insufficient_data" {
		auditDetail = "Analysis aborted: insufficient data quality"
	}

	auditLogRecord := &ExamAnalysisAuditLog{
		ID:          uuid.New(),
		AnalysisID:  &analysisID,
		ActionType:  "process",
		PerformedBy: "SYSTEM_AI_ENGINEER",
		IPAddress:   nil,
		Details:     &auditDetail,
		CreatedAt:   time.Now(),
	}
	_ = worker.repository.CreateAuditLog(ctx, auditLogRecord)

	if worker.eventBus != nil && analysisRecord.Status == "completed" {
		examType := ""
		if analysisRecord.ExamType != nil {
			examType = *analysisRecord.ExamType
		}
		worker.eventBus.Publish(ctx, eventbus.Event{
			Name: "exam.complete",
			Data: map[string]any{
				"title":         "Análise de Exame Concluída",
				"body":          "Laudo disponível para revisão: " + examType,
				"resource_type": "exam_analysis",
				"resource_id":   analysisID.String(),
			},
		})
	}
}

func (worker *Worker) executeAutoCleanup(ctx context.Context) {
	analysesList, err := worker.repository.ListAnalyses(ctx, nil)
	if err != nil {
		slog.Error("Auto-cleanup failed to list analyses", "error", err)
		return
	}

	currentTime := time.Now()
	retentionThreshold := 15 * time.Minute

	for _, analysis := range analysesList {
		if analysis.FilePath == "deleted" || analysis.Status == "pending" || analysis.Status == "processing" {
			continue
		}

		timeElapsed := currentTime.Sub(analysis.CreatedAt)
		if timeElapsed > retentionThreshold {
			if _, statErr := os.Stat(analysis.FilePath); statErr == nil {
				if removeErr := os.Remove(analysis.FilePath); removeErr != nil {
					slog.Error("Failed to physically delete temporary exam file", "filePath", analysis.FilePath, "error", removeErr)
					continue
				}
				slog.Info("Physically deleted temporary medical exam file", "filePath", analysis.FilePath)
			}

			analysis.FilePath = "deleted"
			analysis.UpdatedAt = time.Now()
			_ = worker.repository.UpdateAnalysis(ctx, analysis)

			auditMessage := "Physical temporary file automatically removed due to retention security policy"
			auditRecord := &ExamAnalysisAuditLog{
				ID:          uuid.New(),
				AnalysisID:  &analysis.ID,
				ActionType:  "delete",
				PerformedBy: "SYSTEM_SECURITY_AGENT",
				IPAddress:   nil,
				Details:     &auditMessage,
				CreatedAt:   time.Now(),
			}
			_ = worker.repository.CreateAuditLog(ctx, auditRecord)
		}
	}
}
