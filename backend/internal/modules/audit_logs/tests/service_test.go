package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/healthcare/backend/internal/modules/audit_logs"
	"github.com/healthcare/backend/internal/modules/audit_logs/mocks"
	"github.com/stretchr/testify/assert"
)

func TestService_CreateAuditLog_Success(testingInstance *testing.T) {
	mockRepository := mocks.NewMockAuditRepository()
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
	mockRepository := mocks.NewMockAuditRepository()
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
	mockRepository := mocks.NewMockAuditRepository()
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

func TestService_CreateResourceAuditLog_PersistsRichDetails(testingInstance *testing.T) {
	mockRepository := mocks.NewMockAuditRepository()
	auditService := audit_logs.NewService(mockRepository)
	contextParam := context.Background()

	auditLog, createError := auditService.CreateResourceAuditLog(contextParam, audit_logs.ResourceAuditLog{
		CorrelationID: "correlation-123",
		CallerUserID:  "user-456",
		CallerRole:    "doctor",
		Method:        "UpdateDiagnosticReport",
		AccessGranted: true,
		ResourceType:  "diagnostic_report",
		ResourceID:    "report-789",
		Action:        "update",
		PayloadDiff: map[string]any{
			"conclusion": map[string]any{"before": "old", "after": "new"},
		},
	})

	assert.NoError(testingInstance, createError)
	assert.NotNil(testingInstance, auditLog)
	assert.Equal(testingInstance, "diagnostic_report", auditLog.ResourceType)
	assert.Equal(testingInstance, "report-789", auditLog.ResourceID)
	assert.Equal(testingInstance, "update", auditLog.Action)
	assert.Equal(testingInstance, "correlation-123", auditLog.CorrelationID)
	assert.Len(testingInstance, mockRepository.Logs, 1)
	storedLog := mockRepository.Logs[0]
	assert.Equal(testingInstance, "diagnostic_report", storedLog.ResourceType)
	assert.Equal(testingInstance, map[string]any{
		"conclusion": map[string]any{"before": "old", "after": "new"},
	}, storedLog.PayloadDiff)
}

func TestService_CreateResourceAuditLog_RepositoryFailure(testingInstance *testing.T) {
	mockRepository := mocks.NewMockAuditRepository()
	mockRepository.MockError = errors.New("database insert error")
	auditService := audit_logs.NewService(mockRepository)
	contextParam := context.Background()

	auditLog, createError := auditService.CreateResourceAuditLog(contextParam, audit_logs.ResourceAuditLog{
		CorrelationID: "correlation-123",
		CallerUserID:  "user-456",
		CallerRole:    "doctor",
		Method:        "CancelAppointment",
		AccessGranted: true,
		ResourceType:  "appointment",
		ResourceID:    "appointment-001",
		Action:        "cancel",
	})

	assert.Error(testingInstance, createError)
	assert.Nil(testingInstance, auditLog)
}
