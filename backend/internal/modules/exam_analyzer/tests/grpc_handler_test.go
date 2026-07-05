package tests

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/exam_analyzer"
	"github.com/healthcare/backend/internal/modules/exam_analyzer/mocks"
	"github.com/healthcare/backend/internal/modules/exam_analyzer/pb"
	"github.com/stretchr/testify/assert"
)

func TestGRPCHandler_ListAnalyses(testingInstance *testing.T) {
	mockRepository := mocks.NewMockExamAnalysisRepository()
	serviceInstance := exam_analyzer.NewService(mockRepository, "", "", "")
	grpcHandler := exam_analyzer.NewGRPCHandler(serviceInstance)

	analysisID := uuid.New()
	patientFhirID := "patient-123"
	examType := "Radiografia"
	mockRepository.Analyses[analysisID] = &exam_analyzer.ExamAnalysis{
		ID:            analysisID,
		Status:        "completed",
		PatientFhirID: &patientFhirID,
		ExamType:      &examType,
		FileName:      "rx-chest.png",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	response, err := grpcHandler.ListAnalyses(context.Background(), &pb.ListAnalysesRequest{})

	assert.NoError(testingInstance, err)
	assert.Len(testingInstance, response.Analyses, 1)
	assert.Equal(testingInstance, analysisID.String(), response.Analyses[0].Id)
	assert.Equal(testingInstance, patientFhirID, response.Analyses[0].PatientFhirId)
}

func TestGRPCHandler_ListAnalyses_WithPatientFilter(testingInstance *testing.T) {
	mockRepository := mocks.NewMockExamAnalysisRepository()
	serviceInstance := exam_analyzer.NewService(mockRepository, "", "", "")
	grpcHandler := exam_analyzer.NewGRPCHandler(serviceInstance)

	patientFhirID := "patient-123"
	otherPatientID := "patient-456"
	examType := "Radiografia"

	mockRepository.Analyses[uuid.New()] = &exam_analyzer.ExamAnalysis{
		ID:            uuid.New(),
		Status:        "completed",
		PatientFhirID: &patientFhirID,
		ExamType:      &examType,
		FileName:      "rx-chest.png",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	mockRepository.Analyses[uuid.New()] = &exam_analyzer.ExamAnalysis{
		ID:            uuid.New(),
		Status:        "completed",
		PatientFhirID: &otherPatientID,
		ExamType:      &examType,
		FileName:      "rx-chest-2.png",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	response, err := grpcHandler.ListAnalyses(context.Background(), &pb.ListAnalysesRequest{
		PatientFhirId: patientFhirID,
	})

	assert.NoError(testingInstance, err)
	assert.Len(testingInstance, response.Analyses, 1)
	assert.Equal(testingInstance, patientFhirID, response.Analyses[0].PatientFhirId)
}

func TestGRPCHandler_GetAnalysis(testingInstance *testing.T) {
	mockRepository := mocks.NewMockExamAnalysisRepository()
	serviceInstance := exam_analyzer.NewService(mockRepository, "", "", "")
	grpcHandler := exam_analyzer.NewGRPCHandler(serviceInstance)

	analysisID := uuid.New()
	examType := "Laboratorial"
	mockRepository.Analyses[analysisID] = &exam_analyzer.ExamAnalysis{
		ID:        analysisID,
		Status:    "completed",
		ExamType:  &examType,
		FileName:  "blood-test.pdf",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	response, err := grpcHandler.GetAnalysis(context.Background(), &pb.GetAnalysisRequest{
		AnalysisId: analysisID.String(),
	})

	assert.NoError(testingInstance, err)
	assert.NotNil(testingInstance, response.Analysis)
	assert.Equal(testingInstance, analysisID.String(), response.Analysis.Id)
	assert.Equal(testingInstance, "completed", response.Analysis.Status)
}

func TestGRPCHandler_GetAnalysis_NotFound(testingInstance *testing.T) {
	mockRepository := mocks.NewMockExamAnalysisRepository()
	serviceInstance := exam_analyzer.NewService(mockRepository, "", "", "")
	grpcHandler := exam_analyzer.NewGRPCHandler(serviceInstance)

	response, err := grpcHandler.GetAnalysis(context.Background(), &pb.GetAnalysisRequest{
		AnalysisId: uuid.New().String(),
	})

	assert.Error(testingInstance, err)
	assert.Nil(testingInstance, response)
}

func TestGRPCHandler_GetAnalysis_InvalidID(testingInstance *testing.T) {
	mockRepository := mocks.NewMockExamAnalysisRepository()
	serviceInstance := exam_analyzer.NewService(mockRepository, "", "", "")
	grpcHandler := exam_analyzer.NewGRPCHandler(serviceInstance)

	response, err := grpcHandler.GetAnalysis(context.Background(), &pb.GetAnalysisRequest{
		AnalysisId: "not-a-uuid",
	})

	assert.Error(testingInstance, err)
	assert.Nil(testingInstance, response)
}

func TestGRPCHandler_DeleteAnalysis(testingInstance *testing.T) {
	mockRepository := mocks.NewMockExamAnalysisRepository()
	serviceInstance := exam_analyzer.NewService(mockRepository, "", "", "")
	grpcHandler := exam_analyzer.NewGRPCHandler(serviceInstance)

	analysisID := uuid.New()
	mockRepository.Analyses[analysisID] = &exam_analyzer.ExamAnalysis{
		ID:        analysisID,
		Status:    "completed",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	response, err := grpcHandler.DeleteAnalysis(context.Background(), &pb.DeleteAnalysisRequest{
		AnalysisId: analysisID.String(),
	})

	assert.NoError(testingInstance, err)
	assert.True(testingInstance, response.Success)

	_, exists := mockRepository.Analyses[analysisID]
	assert.False(testingInstance, exists)
}

func TestGRPCHandler_DeleteAnalysis_InvalidID(testingInstance *testing.T) {
	mockRepository := mocks.NewMockExamAnalysisRepository()
	serviceInstance := exam_analyzer.NewService(mockRepository, "", "", "")
	grpcHandler := exam_analyzer.NewGRPCHandler(serviceInstance)

	response, err := grpcHandler.DeleteAnalysis(context.Background(), &pb.DeleteAnalysisRequest{
		AnalysisId: "not-a-uuid",
	})

	assert.Error(testingInstance, err)
	assert.Nil(testingInstance, response)
}
