package audit_logs

import (
	"context"

	"github.com/healthcare/backend/internal/modules/audit_logs/pb"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

type GRPCHandler struct {
	service Service
}

func NewGRPCHandler(service Service) *GRPCHandler {
	return &GRPCHandler{service: service}
}

func (auditLogsHandler *GRPCHandler) CreateAuditLog(contextVal context.Context, request *pb.CreateAuditLogRequest) (*pb.CreateAuditLogResponse, error) {
	violations := make(map[string]string)
	if request.CorrelationId == "" {
		violations["correlation_id"] = "correlation ID is required"
	}
	if request.Method == "" {
		violations["method"] = "method is required"
	}
	if len(violations) > 0 {
		return nil, apperrors.ErrBadRequest.WithFields(violations)
	}

	auditLog, createError := auditLogsHandler.service.CreateAuditLog(
		contextVal,
		request.CorrelationId,
		request.CallerUserId,
		request.CallerRole,
		request.Method,
		request.AccessGranted,
	)
	if createError != nil {
		return nil, apperrors.ToGRPCStatus(createError)
	}

	return &pb.CreateAuditLogResponse{
		Id: auditLog.ID.String(),
	}, nil
}

func (auditLogsHandler *GRPCHandler) ListAuditLogs(contextVal context.Context, request *pb.ListAuditLogsRequest) (*pb.ListAuditLogsResponse, error) {
	logs, totalCount, listError := auditLogsHandler.service.ListAuditLogs(contextVal, int(request.Limit), int(request.Offset))
	if listError != nil {
		return nil, apperrors.ToGRPCStatus(listError)
	}

	pbLogs := make([]*pb.AuditLog, 0, len(logs))
	for _, auditLog := range logs {
		pbLogs = append(pbLogs, &pb.AuditLog{
			Id:            auditLog.ID.String(),
			CorrelationId: auditLog.CorrelationID,
			CallerUserId:  auditLog.CallerUserID,
			CallerRole:    auditLog.CallerRole,
			Method:        auditLog.Method,
			AccessGranted: auditLog.AccessGranted,
			CreatedAt:     auditLog.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &pb.ListAuditLogsResponse{
		AuditLogs: pbLogs,
		Total:     int32(totalCount),
	}, nil
}
