package diagnostic_report

import (
	"context"

	pb "github.com/healthcare/backend/internal/modules/diagnostic_report/pb"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

type GRPCHandler struct {
	service Service
}

func NewGRPCHandler(service Service) *GRPCHandler {
	return &GRPCHandler{service: service}
}

func (handler *GRPCHandler) CreateDiagnosticReport(ctx context.Context, req *pb.CreateDiagnosticReportRequest) (*pb.CreateDiagnosticReportResponse, error) {
	createdReport, err := handler.service.CreateDiagnosticReport(ctx, CreateDiagnosticReportInput{
		EncounterFHIRID: req.EncounterFhirId,
		PatientFHIRID:   req.PatientFhirId,
		ReportCode:      req.ReportCode,
		ReportDisplay:   req.ReportDisplay,
		Conclusion:      req.Conclusion,
	})
	if err != nil {
		return nil, apperrors.ToGRPCStatus(err)
	}

	return &pb.CreateDiagnosticReportResponse{DiagnosticReportFhirId: createdReport.FHIRResourceID}, nil
}

func (handler *GRPCHandler) GetDiagnosticReports(ctx context.Context, req *pb.GetDiagnosticReportsRequest) (*pb.GetDiagnosticReportsResponse, error) {
	if req.EncounterFhirId == "" {
		return nil, apperrors.ErrBadRequest.WithFields(map[string]string{"encounter_fhir_id": "is required"})
	}

	reports, err := handler.service.GetDiagnosticReportsByEncounter(ctx, req.EncounterFhirId)
	if err != nil {
		return nil, apperrors.ToGRPCStatus(err)
	}

	pbReports := make([]*pb.DiagnosticReport, 0, len(reports))
	for _, report := range reports {
		pbReports = append(pbReports, &pb.DiagnosticReport{
			FhirId:        report.FHIRResourceID,
			ReportDisplay: report.ReportDisplay,
			Status:        report.Status,
			Conclusion:    report.Conclusion,
		})
	}

	return &pb.GetDiagnosticReportsResponse{DiagnosticReports: pbReports}, nil
}
