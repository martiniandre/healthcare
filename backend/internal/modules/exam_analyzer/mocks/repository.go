package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/exam_analyzer"
)

type MockExamAnalysisRepository struct {
	Analyses  map[uuid.UUID]*exam_analyzer.ExamAnalysis
	AuditLogs []*exam_analyzer.ExamAnalysisAuditLog
	MockError error
}

func NewMockExamAnalysisRepository() *MockExamAnalysisRepository {
	return &MockExamAnalysisRepository{
		Analyses:  make(map[uuid.UUID]*exam_analyzer.ExamAnalysis),
		AuditLogs: make([]*exam_analyzer.ExamAnalysisAuditLog, 0),
	}
}

func (mockRepo *MockExamAnalysisRepository) CreateAnalysis(ctx context.Context, analysis *exam_analyzer.ExamAnalysis) error {
	if mockRepo.MockError != nil {
		return mockRepo.MockError
	}
	mockRepo.Analyses[analysis.ID] = analysis
	return nil
}

func (mockRepo *MockExamAnalysisRepository) GetAnalysis(ctx context.Context, id uuid.UUID) (*exam_analyzer.ExamAnalysis, error) {
	if mockRepo.MockError != nil {
		return nil, mockRepo.MockError
	}
	analysis, exists := mockRepo.Analyses[id]
	if !exists {
		return nil, exam_analyzer.ErrAnalysisNotFound
	}
	return analysis, nil
}

func (mockRepo *MockExamAnalysisRepository) ListAnalyses(ctx context.Context, patientFhirID *string) ([]*exam_analyzer.ExamAnalysis, error) {
	if mockRepo.MockError != nil {
		return nil, mockRepo.MockError
	}
	result := make([]*exam_analyzer.ExamAnalysis, 0, len(mockRepo.Analyses))
	for _, analysis := range mockRepo.Analyses {
		if patientFhirID != nil && *patientFhirID != "" && (analysis.PatientFhirID == nil || *analysis.PatientFhirID != *patientFhirID) {
			continue
		}
		result = append(result, analysis)
	}
	return result, nil
}

func (mockRepo *MockExamAnalysisRepository) UpdateAnalysis(ctx context.Context, analysis *exam_analyzer.ExamAnalysis) error {
	if mockRepo.MockError != nil {
		return mockRepo.MockError
	}
	mockRepo.Analyses[analysis.ID] = analysis
	return nil
}

func (mockRepo *MockExamAnalysisRepository) DeleteAnalysis(ctx context.Context, id uuid.UUID) error {
	if mockRepo.MockError != nil {
		return mockRepo.MockError
	}
	delete(mockRepo.Analyses, id)
	return nil
}

func (mockRepo *MockExamAnalysisRepository) CreateAuditLog(ctx context.Context, log *exam_analyzer.ExamAnalysisAuditLog) error {
	if mockRepo.MockError != nil {
		return mockRepo.MockError
	}
	mockRepo.AuditLogs = append(mockRepo.AuditLogs, log)
	return nil
}
