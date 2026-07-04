package tests

import (
	"context"
	"testing"

	"github.com/healthcare/backend/internal/modules/diagnostic_report"
	"github.com/healthcare/backend/internal/modules/diagnostic_report/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateDiagnosticReport_DefaultsToFinal(t *testing.T) {
	diagnosticReportService := diagnostic_report.NewService(&mocks.MockDiagnosticReportRepository{})

	entity := &diagnostic_report.DiagnosticReport{
		PatientFHIRID:   "patient-123",
		EncounterFHIRID: "encounter-456",
		ReportCode:      "24323-8",
		ReportDisplay:   "Complete blood count",
		Conclusion:      "Normal values",
	}

	result, err := diagnosticReportService.CreateDiagnosticReport(context.Background(), entity)

	assert.NoError(t, err)
	assert.Equal(t, "final", result.Status)
}
