package diagnostic_report

import (
	"context"
	"errors"

	pb "github.com/healthcare/backend/internal/modules/diagnostic_report/pb"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

type GRPCHandler struct {
	service Service
}

func NewGRPCHandler(service Service) *GRPCHandler {
	return &GRPCHandler{service: service}
}

func mapDiagnosticReportError(err error) error {
	if errors.Is(err, ErrDiagnosticReportNotFound) {
		return apperrors.ErrDiagnosticReportNotFound.ToGRPC()
	}
	return apperrors.ToGRPCStatus(err)
}

func (handler *GRPCHandler) CreateDiagnosticReport(ctx context.Context, req *pb.CreateDiagnosticReportRequest) (*pb.CreateDiagnosticReportResponse, error) {
	violations := make(map[string]string)
	if req.PatientFhirId == "" {
		violations["patient_fhir_id"] = "is required"
	}
	if req.ReportCode == "" {
		violations["report_code"] = "is required"
	}
	if req.EncounterFhirId == "" {
		violations["encounter_fhir_id"] = "is required"
	}
	if len(violations) > 0 {
		return nil, apperrors.ErrBadRequest.WithFields(violations)
	}

	report := &DiagnosticReport{
		EncounterFHIRID: req.EncounterFhirId,
		PatientFHIRID:   req.PatientFhirId,
		ReportCode:      req.ReportCode,
		ReportDisplay:   req.ReportDisplay,
		Conclusion:      req.Conclusion,
	}

	createdReport, err := handler.service.CreateDiagnosticReport(ctx, report)
	if err != nil {
		return nil, mapDiagnosticReportError(err)
	}

	return &pb.CreateDiagnosticReportResponse{DiagnosticReportFhirId: createdReport.FHIRResourceID}, nil
}

func (handler *GRPCHandler) GetDiagnosticReports(ctx context.Context, req *pb.GetDiagnosticReportsRequest) (*pb.GetDiagnosticReportsResponse, error) {
	if req.EncounterFhirId == "" {
		return nil, apperrors.ErrBadRequest.WithFields(map[string]string{"encounter_fhir_id": "is required"})
	}

	reports, err := handler.service.GetDiagnosticReportsByEncounter(ctx, req.EncounterFhirId)
	if err != nil {
		return nil, mapDiagnosticReportError(err)
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
