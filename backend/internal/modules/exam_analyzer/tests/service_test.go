package tests

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/healthcare/backend/internal/modules/exam_analyzer"
	"github.com/healthcare/backend/internal/modules/exam_analyzer/mocks"
	"github.com/stretchr/testify/assert"
)

func stringsContainsProbabilisticLanguage(text string) bool {
	lowerText := strings.ToLower(text)
	probabilisticTerms := []string{"pode", "sugerir", "compatibilidade", "compatível", "achado"}
	for _, term := range probabilisticTerms {
		if strings.Contains(lowerText, term) {
			return true
		}
	}
	return false
}

func TestService_AnalyzeExamFile_InsufficientData(testingInstance *testing.T) {
	mockRepository := mocks.NewMockExamAnalysisRepository()
	serviceInstance := exam_analyzer.NewService(mockRepository, "", "", "")
	contextParam := context.Background()

	temporaryDirectory := testingInstance.TempDir()
	temporaryFilePath := filepath.Join(temporaryDirectory, "low_res_exam.png")

	writeError := os.WriteFile(temporaryFilePath, make([]byte, 100), 0644)
	assert.NoError(testingInstance, writeError)

	analysisResponse, statusResult, err := serviceInstance.AnalyzeExamFile(contextParam, temporaryFilePath, "low_res_exam.png")

	assert.NoError(testingInstance, err)
	assert.Equal(testingInstance, "insufficient_data", statusResult)
	assert.Nil(testingInstance, analysisResponse)
}

func TestService_AnalyzeExamFile_Simulation_Xray(testingInstance *testing.T) {
	mockRepository := mocks.NewMockExamAnalysisRepository()
	serviceInstance := exam_analyzer.NewService(mockRepository, "", "", "")
	contextParam := context.Background()

	temporaryDirectory := testingInstance.TempDir()
	temporaryFilePath := filepath.Join(temporaryDirectory, "rx_chest.png")

	writeError := os.WriteFile(temporaryFilePath, make([]byte, 6000), 0644)
	assert.NoError(testingInstance, writeError)

	analysisResponse, statusResult, err := serviceInstance.AnalyzeExamFile(contextParam, temporaryFilePath, "rx_chest.png")

	assert.NoError(testingInstance, err)
	assert.Equal(testingInstance, "completed", statusResult)
	assert.NotNil(testingInstance, analysisResponse)
	assert.Contains(testingInstance, analysisResponse.ExamType, "Radiografia")

	foundProbabilisticWord := false
	for _, finding := range analysisResponse.DetectedFindings {
		if stringsContainsProbabilisticLanguage(finding.Finding) {
			foundProbabilisticWord = true
			break
		}
	}
	assert.True(testingInstance, foundProbabilisticWord, "Should contain probabilistic language in clinical findings")
}

func TestService_AnalyzeExamFile_Simulation_Pdf(testingInstance *testing.T) {
	mockRepository := mocks.NewMockExamAnalysisRepository()
	serviceInstance := exam_analyzer.NewService(mockRepository, "", "", "")
	contextParam := context.Background()

	temporaryDirectory := testingInstance.TempDir()
	temporaryFilePath := filepath.Join(temporaryDirectory, "blood_test.pdf")

	writeError := os.WriteFile(temporaryFilePath, make([]byte, 6000), 0644)
	assert.NoError(testingInstance, writeError)

	analysisResponse, statusResult, err := serviceInstance.AnalyzeExamFile(contextParam, temporaryFilePath, "blood_test.pdf")

	assert.NoError(testingInstance, err)
	assert.Equal(testingInstance, "completed", statusResult)
	assert.NotNil(testingInstance, analysisResponse)
	assert.Contains(testingInstance, analysisResponse.ExamType, "Laboratorial")
}
