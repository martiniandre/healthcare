package diagnostic_report

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/fhir"
	"github.com/healthcare/backend/internal/shared/healthcare"
)

type Repository interface {
	CreateDiagnosticReport(ctx context.Context, report *DiagnosticReport) (*DiagnosticReport, error)
	GetDiagnosticReportByID(ctx context.Context, reportFHIRID string) (*DiagnosticReport, error)
	UpdateDiagnosticReport(ctx context.Context, reportFHIRID string, report *DiagnosticReport) (*DiagnosticReport, error)
	GetDiagnosticReportsByEncounter(ctx context.Context, encounterFHIRID string) ([]*DiagnosticReport, error)
}

type repository struct {
	fhirClient healthcare.FHIRClient
}

func NewRepository(fhirClient healthcare.FHIRClient) Repository {
	return &repository{fhirClient: fhirClient}
}

func (reportRepository *repository) CreateDiagnosticReport(ctx context.Context, report *DiagnosticReport) (*DiagnosticReport, error) {
	fhirReport := fhir.NewDiagnosticReportResource(
		report.PatientFHIRID,
		report.EncounterFHIRID,
		report.ReportCode,
		report.ReportDisplay,
		report.Conclusion,
	)

	responseBody, err := reportRepository.fhirClient.CreateResource(ctx, "DiagnosticReport", fhirReport)
	if err != nil {
		if healthcare.IsNotFound(err) {
			return nil, apperrors.ErrDiagnosticReportNotFound
		}
		return nil, fmt.Errorf("failed to create diagnostic report: %w", err)
	}

	decodedResource, err := fhir.DecodeResource[fhir.DiagnosticReport](responseBody)
	if err != nil {
		return nil, fmt.Errorf("failed to parse diagnostic report response: %w", err)
	}

	report.FHIRResourceID = decodedResource.ID
	report.Version = extractVersionID(decodedResource.Meta)
	return report, nil
}

func (reportRepository *repository) GetDiagnosticReportByID(ctx context.Context, reportFHIRID string) (*DiagnosticReport, error) {
	responseBody, err := reportRepository.fhirClient.GetResource(ctx, "DiagnosticReport", reportFHIRID)
	if err != nil {
		if healthcare.IsNotFound(err) {
			return nil, apperrors.ErrDiagnosticReportNotFound
		}
		return nil, fmt.Errorf("failed to get diagnostic report: %w", err)
	}

	decodedResource, err := fhir.DecodeResource[fhir.DiagnosticReport](responseBody)
	if err != nil {
		return nil, fmt.Errorf("failed to parse diagnostic report response: %w", err)
	}
	return mapFHIRDiagnosticReportToDomain(decodedResource), nil
}

func (reportRepository *repository) UpdateDiagnosticReport(ctx context.Context, reportFHIRID string, report *DiagnosticReport) (*DiagnosticReport, error) {
	fhirReport := &fhir.DiagnosticReportResource{
		ResourceType: "DiagnosticReport",
		ID:           reportFHIRID,
		Status:       report.Status,
		Code: fhir.CodeableConcept{
			Coding: []fhir.Coding{
				{System: "http://loinc.org", Code: report.ReportCode, Display: report.ReportDisplay},
			},
			Text: report.ReportDisplay,
		},
		Subject:    fhir.Reference{Reference: "Patient/" + report.PatientFHIRID},
		Encounter:  fhir.Reference{Reference: "Encounter/" + report.EncounterFHIRID},
		Issued:     report.IssuedAt.Format(time.RFC3339),
		Conclusion: report.Conclusion,
	}

	responseBody, err := reportRepository.fhirClient.UpdateResource(ctx, "DiagnosticReport", reportFHIRID, fhirReport)
	if err != nil {
		if healthcare.IsNotFound(err) {
			return nil, apperrors.ErrDiagnosticReportNotFound
		}
		return nil, fmt.Errorf("failed to update diagnostic report: %w", err)
	}

	decodedResource, err := fhir.DecodeResource[fhir.DiagnosticReport](responseBody)
	if err != nil {
		return nil, fmt.Errorf("failed to parse diagnostic report response: %w", err)
	}
	updatedReport := mapFHIRDiagnosticReportToDomain(decodedResource)
	if updatedReport.FHIRResourceID == "" {
		updatedReport.FHIRResourceID = reportFHIRID
	}
	return updatedReport, nil
}

func extractVersionID(meta *fhir.ResourceMeta) string {
	if meta != nil && meta.VersionID != "" {
		return meta.VersionID
	}
	return "1"
}

func (reportRepository *repository) GetDiagnosticReportsByEncounter(ctx context.Context, encounterFHIRID string) ([]*DiagnosticReport, error) {
	queryParams := url.Values{"encounter": []string{fmt.Sprintf("Encounter/%s", encounterFHIRID)}}.Encode()
	responseBody, err := reportRepository.fhirClient.SearchResources(ctx, "DiagnosticReport", queryParams)
	if err != nil {
		if healthcare.IsNotFound(err) {
			return nil, apperrors.ErrDiagnosticReportNotFound
		}
		return nil, fmt.Errorf("failed to search diagnostic reports: %w", err)
	}
	return parseDiagnosticReportBundle(responseBody)
}

func parseDiagnosticReportBundle(responseBody json.RawMessage) ([]*DiagnosticReport, error) {
	decodedResources, err := fhir.DecodeBundle[fhir.DiagnosticReport](responseBody)
	if err != nil {
		return nil, err
	}
	reports := make([]*DiagnosticReport, 0, len(decodedResources))
	for _, resource := range decodedResources {
		reports = append(reports, mapFHIRDiagnosticReportToDomain(&resource))
	}
	return reports, nil
}

func mapFHIRDiagnosticReportToDomain(resource *fhir.DiagnosticReport) *DiagnosticReport {
	report := &DiagnosticReport{}
	report.FHIRResourceID = resource.ID
	report.Status = resource.Status
	report.Conclusion = resource.Conclusion
	code, display, text := fhir.CodeableConceptParts(resource.Code)
	report.ReportDisplay = text
	report.ReportCode = code
	if report.ReportDisplay == "" {
		report.ReportDisplay = display
	}
	report.EncounterFHIRID = fhir.SplitReferenceID(resource.Encounter.Reference)
	report.PatientFHIRID = fhir.SplitReferenceID(resource.Subject.Reference)
	if parsedTime, ok := fhir.ParseRFC3339(resource.Issued); ok {
		report.IssuedAt = parsedTime
	}
	report.Version = extractVersionID(resource.Meta)
	return report
}
