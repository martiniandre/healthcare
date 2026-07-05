package exam_analyzer

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/exam_analyzer/pb"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

func mapExamAnalyzerError(err error) error {
	if errors.Is(err, ErrAnalysisNotFound) {
		return apperrors.ErrNotFound.ToGRPC()
	}
	return apperrors.ToGRPCStatus(err)
}

type GRPCHandler struct {
	service Service
	repo    Repository
}

func NewGRPCHandler(repo Repository, service Service) *GRPCHandler {
	return &GRPCHandler{service: service, repo: repo}
}

func (handler *GRPCHandler) ListAnalyses(ctx context.Context, req *pb.ListAnalysesRequest) (*pb.ListAnalysesResponse, error) {
	var filterPatient *string
	if req.PatientFhirId != "" {
		filterPatient = &req.PatientFhirId
	}

	analyses, err := handler.repo.ListAnalyses(ctx, filterPatient)
	if err != nil {
		return nil, mapExamAnalyzerError(err)
	}

	items := make([]*pb.ExamAnalysisItem, 0, len(analyses))
	for _, a := range analyses {
		userID := ""
		if a.UserID != nil {
			userID = a.UserID.String()
		}
		patientFhirID := ""
		if a.PatientFhirID != nil {
			patientFhirID = *a.PatientFhirID
		}
		examType := ""
		if a.ExamType != nil {
			examType = *a.ExamType
		}

		items = append(items, &pb.ExamAnalysisItem{
			Id:               a.ID.String(),
			UserId:           userID,
			PatientFhirId:    patientFhirID,
			ExamType:         examType,
			FileName:         a.FileName,
			Status:           a.Status,
			AnalysisResponse: string(a.AnalysisResponse),
			ConsentGiven:     a.ConsentGiven,
			Anonymized:       a.Anonymized,
			CreatedAt:        a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:        a.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &pb.ListAnalysesResponse{Analyses: items}, nil
}

func (handler *GRPCHandler) GetAnalysis(ctx context.Context, req *pb.GetAnalysisRequest) (*pb.GetAnalysisResponse, error) {
	analysisID, err := uuid.Parse(req.AnalysisId)
	if err != nil {
		return nil, apperrors.ErrBadRequest.ToGRPC()
	}

	analysis, err := handler.repo.GetAnalysis(ctx, analysisID)
	if err != nil {
		return nil, mapExamAnalyzerError(err)
	}

	userID := ""
	if analysis.UserID != nil {
		userID = analysis.UserID.String()
	}
	patientFhirID := ""
	if analysis.PatientFhirID != nil {
		patientFhirID = *analysis.PatientFhirID
	}
	examType := ""
	if analysis.ExamType != nil {
		examType = *analysis.ExamType
	}

	return &pb.GetAnalysisResponse{
		Analysis: &pb.ExamAnalysisItem{
			Id:               analysis.ID.String(),
			UserId:           userID,
			PatientFhirId:    patientFhirID,
			ExamType:         examType,
			FileName:         analysis.FileName,
			Status:           analysis.Status,
			AnalysisResponse: string(analysis.AnalysisResponse),
			ConsentGiven:     analysis.ConsentGiven,
			Anonymized:       analysis.Anonymized,
			CreatedAt:        analysis.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:        analysis.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

func (handler *GRPCHandler) DeleteAnalysis(ctx context.Context, req *pb.DeleteAnalysisRequest) (*pb.DeleteAnalysisResponse, error) {
	analysisID, err := uuid.Parse(req.AnalysisId)
	if err != nil {
		return nil, apperrors.ErrBadRequest.ToGRPC()
	}

	err = handler.repo.DeleteAnalysis(ctx, analysisID)
	if err != nil {
		return nil, mapExamAnalyzerError(err)
	}

	return &pb.DeleteAnalysisResponse{Success: true}, nil
}
