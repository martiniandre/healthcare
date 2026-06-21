package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/healthcare/backend/internal/modules/audit_logs"
	"github.com/stretchr/testify/assert"
)

type MockAuditRepository struct {
	Logs      []*audit_logs.AuditLog
	MockError error
}

func NewMockAuditRepository() *MockAuditRepository {
	return &MockAuditRepository{
		Logs: make([]*audit_logs.AuditLog, 0),
	}
}

func (mockRepository *MockAuditRepository) CreateAuditLog(contextVal context.Context, auditLog *audit_logs.AuditLog) error {
	if mockRepository.MockError != nil {
		return mockRepository.MockError
	}
	mockRepository.Logs = append(mockRepository.Logs, auditLog)
	return nil
}

func (mockRepository *MockAuditRepository) ListAuditLogs(contextVal context.Context, limit int, offset int) ([]*audit_logs.AuditLog, int, error) {
	if mockRepository.MockError != nil {
		return nil, 0, mockRepository.MockError
	}

	totalCount := len(mockRepository.Logs)
	if offset >= totalCount {
		return []*audit_logs.AuditLog{}, totalCount, nil
	}

	endIndex := offset + limit
	if endIndex > totalCount {
		endIndex = totalCount
	}

	return mockRepository.Logs[offset:endIndex], totalCount, nil
}

func TestService_CreateAuditLog_Success(testingInstance *testing.T) {
	mockRepository := NewMockAuditRepository()
	auditService := audit_logs.NewService(mockRepository)
	contextParam := context.Background()

	auditLog, createError := auditService.CreateAuditLog(
		contextParam,
		"correlation-123",
		"user-456",
		"doctor",
		"/clinical.v1.ClinicalService/CreateEncounter",
		true,
	)

	assert.NoError(testingInstance, createError)
	assert.NotNil(testingInstance, auditLog)
	assert.Equal(testingInstance, "correlation-123", auditLog.CorrelationID)
	assert.Equal(testingInstance, "user-456", auditLog.CallerUserID)
	assert.Equal(testingInstance, "doctor", auditLog.CallerRole)
	assert.Equal(testingInstance, "/clinical.v1.ClinicalService/CreateEncounter", auditLog.Method)
	assert.True(testingInstance, auditLog.AccessGranted)
	assert.Len(testingInstance, mockRepository.Logs, 1)
}

func TestService_CreateAuditLog_Failure(testingInstance *testing.T) {
	mockRepository := NewMockAuditRepository()
	mockRepository.MockError = errors.New("database insert error")
	auditService := audit_logs.NewService(mockRepository)
	contextParam := context.Background()

	auditLog, createError := auditService.CreateAuditLog(
		contextParam,
		"correlation-123",
		"user-456",
		"doctor",
		"/clinical.v1.ClinicalService/CreateEncounter",
		true,
	)

	assert.Error(testingInstance, createError)
	assert.Nil(testingInstance, auditLog)
}

func TestService_ListAuditLogs(testingInstance *testing.T) {
	mockRepository := NewMockAuditRepository()
	auditService := audit_logs.NewService(mockRepository)
	contextParam := context.Background()

	for iterationIndex := 0; iterationIndex < 15; iterationIndex++ {
		_, createError := auditService.CreateAuditLog(
			contextParam,
			"correlation-id",
			"user-id",
			"nurse",
			"some-method",
			true,
		)
		assert.NoError(testingInstance, createError)
	}

	logs, totalCount, listError := auditService.ListAuditLogs(contextParam, 10, 0)
	assert.NoError(testingInstance, listError)
	assert.Equal(testingInstance, 15, totalCount)
	assert.Len(testingInstance, logs, 10)

	logsPageTwo, totalCountPageTwo, listErrorPageTwo := auditService.ListAuditLogs(contextParam, 10, 10)
	assert.NoError(testingInstance, listErrorPageTwo)
	assert.Equal(testingInstance, 15, totalCountPageTwo)
	assert.Len(testingInstance, logsPageTwo, 5)
}
