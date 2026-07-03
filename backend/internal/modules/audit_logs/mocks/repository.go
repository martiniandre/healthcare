package mocks

import (
	"context"

	"github.com/healthcare/backend/internal/modules/audit_logs"
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
