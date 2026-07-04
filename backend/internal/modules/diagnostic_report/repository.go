package diagnostic_report

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/healthcare/backend/internal/shared/fhir"
	"github.com/healthcare/backend/internal/shared/healthcare"
)

type Repository interface {
	CreateDiagnosticReport(ctx context.Context, report *DiagnosticReport) (*DiagnosticReport, error)
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
		return nil, fmt.Errorf("failed to create diagnostic report: %w", err)
	}

	var createdResource map[string]interface{}
	if err := json.Unmarshal(responseBody, &createdResource); err != nil {
		return nil, fmt.Errorf("failed to parse diagnostic report response: %w", err)
	}

	fhirID, _ := createdResource["id"].(string)
	report.FHIRResourceID = fhirID
	report.IssuedAt = time.Now()
	return report, nil
}

func (reportRepository *repository) GetDiagnosticReportsByEncounter(ctx context.Context, encounterFHIRID string) ([]*DiagnosticReport, error) {
	queryParams := url.Values{"encounter": []string{fmt.Sprintf("Encounter/%s", encounterFHIRID)}}.Encode()
	responseBody, err := reportRepository.fhirClient.SearchResources(ctx, "DiagnosticReport", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search diagnostic reports: %w", err)
	}
	return parseDiagnosticReportBundle(responseBody)
}

func extractBundleEntries(responseBody json.RawMessage) ([]map[string]interface{}, error) {
	var bundle map[string]interface{}
	if err := json.Unmarshal(responseBody, &bundle); err != nil {
		return nil, err
	}
	rawEntries, ok := bundle["entry"].([]interface{})
	if !ok {
		return []map[string]interface{}{}, nil
	}
	entries := make([]map[string]interface{}, 0, len(rawEntries))
	for _, rawEntry := range rawEntries {
		entryMap, ok := rawEntry.(map[string]interface{})
		if !ok {
			continue
		}
		resource, ok := entryMap["resource"].(map[string]interface{})
		if !ok {
			continue
		}
		entries = append(entries, resource)
	}
	return entries, nil
}

func parseDiagnosticReportBundle(responseBody json.RawMessage) ([]*DiagnosticReport, error) {
	entries, err := extractBundleEntries(responseBody)
	if err != nil {
		return nil, err
	}
	reports := make([]*DiagnosticReport, 0, len(entries))
	for _, resource := range entries {
		report := &DiagnosticReport{}
		report.FHIRResourceID, _ = resource["id"].(string)
		report.Status, _ = resource["status"].(string)
		report.Conclusion, _ = resource["conclusion"].(string)
		if codes, ok := resource["code"].(map[string]interface{}); ok {
			report.ReportDisplay, _ = codes["text"].(string)
		}
		reports = append(reports, report)
	}
	return reports, nil
}
