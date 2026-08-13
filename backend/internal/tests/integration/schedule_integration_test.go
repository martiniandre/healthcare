package integration

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestAppointmentLifecycleAndOverlapConflict(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	receptionClient := loginAs(t, serverURL, "recepcao@hospital.com", "secret123")

	doctorEmployeeID := seedDoctorEmployee(t, testServer.db)
	tomorrowStartsAt := time.Now().Add(24 * time.Hour).Truncate(time.Second)
	tomorrowEndsAt := tomorrowStartsAt.Add(30 * time.Minute)

	createResponse := receptionClient.Post(t, "/api/v1/appointments", map[string]interface{}{
		"patient_fhir_id": "patient-1",
		"staff_id":        doctorEmployeeID.String(),
		"starts_at":       tomorrowStartsAt.Format(time.RFC3339),
		"ends_at":         tomorrowEndsAt.Format(time.RFC3339),
		"reason":          "Consulta de rotina",
	})
	requireStatusCode(t, createResponse, http.StatusCreated)

	appointmentID := requireJSONFieldValue(t, createResponse, "id")

	getResponse := receptionClient.Get(t, "/api/v1/appointments/"+appointmentID)
	requireStatusCode(t, getResponse, http.StatusOK)
	status := requireJSONFieldValue(t, getResponse, "status")
	if status != "scheduled" {
		t.Fatalf("expected status scheduled, got %q", status)
	}

	overlapResponse := receptionClient.Post(t, "/api/v1/appointments", map[string]interface{}{
		"patient_fhir_id": "patient-2",
		"staff_id":        doctorEmployeeID.String(),
		"starts_at":       tomorrowStartsAt.Add(15 * time.Minute).Format(time.RFC3339),
		"ends_at":         tomorrowEndsAt.Add(15 * time.Minute).Format(time.RFC3339),
		"reason":          "Consulta com conflito",
	})
	requireStatusCode(t, overlapResponse, http.StatusConflict)

	cancelResponse := receptionClient.Post(t, "/api/v1/appointments/"+appointmentID+"/cancel", nil)
	requireStatusCode(t, cancelResponse, http.StatusOK)
	cancelledStatus := requireJSONFieldValue(t, cancelResponse, "status")
	if cancelledStatus != "cancelled" {
		t.Fatalf("expected status cancelled, got %q", cancelledStatus)
	}

	requireStatusCode(t, receptionClient.Get(t, "/api/v1/appointments?patient_fhir_id=patient-1"), http.StatusOK)
	requireStatusCode(t, receptionClient.Get(t, "/api/v1/appointments?staff_id="+doctorEmployeeID.String()+"&date="+tomorrowStartsAt.Format("2006-01-02")), http.StatusOK)
}

func TestCreateAppointmentRequiresFutureStart(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	receptionClient := loginAs(t, serverURL, "recepcao@hospital.com", "secret123")
	doctorEmployeeID := seedDoctorEmployee(t, testServer.db)

	pastStartsAt := time.Now().Add(-2 * time.Hour).Truncate(time.Second)
	pastEndsAt := pastStartsAt.Add(30 * time.Minute)

	response := receptionClient.Post(t, "/api/v1/appointments", map[string]interface{}{
		"patient_fhir_id": "patient-1",
		"staff_id":        doctorEmployeeID.String(),
		"starts_at":       pastStartsAt.Format(time.RFC3339),
		"ends_at":         pastEndsAt.Format(time.RFC3339),
		"reason":          "Consulta no passado",
	})
	requireStatusCode(t, response, http.StatusBadRequest)
}

func TestCreateAppointmentRejectsNonExistentStaff(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	receptionClient := loginAs(t, serverURL, "recepcao@hospital.com", "secret123")
	tomorrowStartsAt := time.Now().Add(24 * time.Hour).Truncate(time.Second)

	response := receptionClient.Post(t, "/api/v1/appointments", map[string]interface{}{
		"patient_fhir_id": "patient-1",
		"staff_id":        uuid.New().String(),
		"starts_at":       tomorrowStartsAt.Format(time.RFC3339),
		"ends_at":         tomorrowStartsAt.Add(30 * time.Minute).Format(time.RFC3339),
		"reason":          "Consulta com staff inexistente",
	})
	requireStatusCode(t, response, http.StatusNotFound)
}

func TestPatientCanListOwnAppointments(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	patientUserID := fetchUserID(t, testServer.db, "paciente@mail.com")

	receptionClient := loginAs(t, serverURL, "recepcao@hospital.com", "secret123")
	doctorEmployeeID := seedDoctorEmployee(t, testServer.db)
	tomorrowStartsAt := time.Now().Add(24 * time.Hour).Truncate(time.Second)

	createResponse := receptionClient.Post(t, "/api/v1/appointments", map[string]interface{}{
		"patient_fhir_id": patientUserID,
		"staff_id":        doctorEmployeeID.String(),
		"starts_at":       tomorrowStartsAt.Format(time.RFC3339),
		"ends_at":         tomorrowStartsAt.Add(30 * time.Minute).Format(time.RFC3339),
		"reason":          "Consulta de rotina",
	})
	requireStatusCode(t, createResponse, http.StatusCreated)

	patientClient := loginAs(t, serverURL, "paciente@mail.com", "secret123")
	myAppointmentsResponse := patientClient.Get(t, "/api/v1/appointments/my")
	requireStatusCode(t, myAppointmentsResponse, http.StatusOK)

	var appointmentsList []map[string]interface{}
	decodeJSONResponse(t, myAppointmentsResponse, &appointmentsList)
	if len(appointmentsList) == 0 {
		t.Fatal("expected at least one appointment for the patient")
	}
}

func fetchUserID(t *testing.T, db *pgxpool.Pool, email string) string {
	t.Helper()
	ctx := context.Background()
	var userID uuid.UUID
	if queryError := db.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, email).Scan(&userID); queryError != nil {
		t.Fatalf("failed to fetch user id for %s: %v", email, queryError)
	}
	return userID.String()
}

func seedDoctorEmployee(t *testing.T, db *pgxpool.Pool) uuid.UUID {
	t.Helper()
	ctx := context.Background()

	var employeeID uuid.UUID
	insertError := db.QueryRow(ctx, `
		INSERT INTO employees (id, full_name, email, role, crm_number, is_active, created_at, updated_at)
		VALUES (uuid_generate_v4(), 'Médico Teste', 'medico.teste@clinica.com', 'DOCTOR', 'CRM 123456', true, NOW(), NOW())
		RETURNING id`).Scan(&employeeID)
	if insertError != nil {
		t.Fatalf("failed to seed doctor employee: %v", insertError)
	}
	return employeeID
}
