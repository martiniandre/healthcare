package integration

import (
	"net/http"
	"testing"
	"time"
)

func TestAuditLogsRecordAppointmentActions(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	receptionClient := loginAs(t, serverURL, "recepcao@hospital.com", "secret123")
	doctorEmployeeID := seedDoctorEmployee(t, testServer.db)
	tomorrowStartsAt := alignToValidSlot(time.Now().Add(24 * time.Hour))

	createResponse := receptionClient.Post(t, "/api/v1/appointments", map[string]interface{}{
		"patient_fhir_id": "patient-1",
		"staff_id":        doctorEmployeeID.String(),
		"starts_at":       tomorrowStartsAt.Format(time.RFC3339),
		"ends_at":         tomorrowStartsAt.Add(30 * time.Minute).Format(time.RFC3339),
		"reason":          "Consulta de rotina",
	})
	requireStatusCode(t, createResponse, http.StatusCreated)

	adminClient := loginAs(t, serverURL, "admin@hospital.com", "secret123")
	auditResponse := adminClient.Get(t, "/api/v1/audit-logs")
	requireStatusCode(t, auditResponse, http.StatusOK)
	var auditBody struct {
		AuditLogs []map[string]interface{} `json:"audit_logs"`
		Total     int                      `json:"total"`
	}
	decodeJSONResponse(t, auditResponse, &auditBody)
	if auditBody.Total == 0 || len(auditBody.AuditLogs) == 0 {
		t.Fatal("expected audit logs after appointment creation")
	}
}

func TestOnlyAdminCanListAuditLogs(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	nurseClient := loginAs(t, serverURL, "enfermeiro@hospital.com", "secret123")
	requireStatusCode(t, nurseClient.Get(t, "/api/v1/audit-logs"), http.StatusForbidden)
}
