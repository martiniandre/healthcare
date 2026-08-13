package integration

import (
	"net/http"
	"testing"
)

func TestPatientRoleCannotAccessStaffPatientEndpoints(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	patientClient := loginAs(t, serverURL, "paciente@mail.com", "secret123")

	requireStatusCode(t, patientClient.Get(t, "/api/v1/patients"), http.StatusForbidden)
	requireStatusCode(t, patientClient.Get(t, "/api/v1/appointments"), http.StatusForbidden)
	requireStatusCode(t, patientClient.Get(t, "/api/v1/analytics/dashboard"), http.StatusForbidden)
	requireStatusCode(t, patientClient.Get(t, "/api/v1/staff/employees"), http.StatusForbidden)
}

func TestPatientRoleCanAccessOwnPortal(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	patientClient := loginAs(t, serverURL, "paciente@mail.com", "secret123")

	requireStatusCode(t, patientClient.Get(t, "/api/v1/portal/dashboard"), http.StatusOK)
	requireStatusCode(t, patientClient.Get(t, "/api/v1/appointments/my"), http.StatusOK)
	requireStatusCode(t, patientClient.Get(t, "/api/v1/notifications/unread-count"), http.StatusOK)
}

func TestDoctorRoleCannotAccessAdminOnlyEndpoints(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	doctorClient := loginAs(t, serverURL, "medico@clinica.com", "secret123")

	requireStatusCode(t, doctorClient.Get(t, "/api/v1/audit-logs"), http.StatusForbidden)
	requireStatusCode(t, doctorClient.Post(t, "/api/v1/staff/employees", map[string]string{
		"full_name": "Novo Funcionário",
		"email":     "novo@hospital.com",
		"role":      "NURSE",
	}), http.StatusForbidden)
}

func TestDoctorRoleCanAccessClinicalEndpoints(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	doctorClient := loginAs(t, serverURL, "medico@clinica.com", "secret123")

	requireStatusCode(t, doctorClient.Get(t, "/api/v1/patients"), http.StatusOK)
	requireStatusCode(t, doctorClient.Get(t, "/api/v1/telemetry/rooms"), http.StatusOK)
	requireStatusCode(t, doctorClient.Get(t, "/api/v1/analytics/dashboard"), http.StatusOK)
}

func TestReceptionRoleCanManagePatientsButNotStaff(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	receptionClient := loginAs(t, serverURL, "recepcao@hospital.com", "secret123")

	requireStatusCode(t, receptionClient.Post(t, "/api/v1/staff/employees", map[string]string{
		"full_name": "Novo Funcionário",
		"email":     "novo@hospital.com",
		"role":      "NURSE",
	}), http.StatusForbidden)

	patientResponse := receptionClient.Post(t, "/api/v1/patients", map[string]string{
		"full_name":    "Maria Silva",
		"birth_date":   "1990-01-01",
		"document_id":  "529.982.247-25",
		"phone_number": "(11) 98765-4321",
	})
	requireStatusCode(t, patientResponse, http.StatusCreated)
}

func TestNurseRoleCannotPrescribeMedication(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	nurseClient := loginAs(t, serverURL, "enfermeiro@hospital.com", "secret123")

	requireStatusCode(t, nurseClient.Post(t, "/api/v1/encounters/encounter-1/medications", map[string]string{
		"patient_fhir_id": "patient-1",
		"medication_name": "Dipirona 500mg",
	}), http.StatusForbidden)
}
