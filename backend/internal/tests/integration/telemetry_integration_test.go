package integration

import (
	"context"
	"net/http"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestTelemetryRoomBedAndUnlockFlow(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	roomID, bedID := seedTelemetryRoomAndBed(t, testServer.db)

	nurseClient := loginAs(t, serverURL, "enfermeiro@hospital.com", "secret123")

	roomsResponse := nurseClient.Get(t, "/api/v1/telemetry/rooms")
	requireStatusCode(t, roomsResponse, http.StatusOK)
	var roomsList []map[string]interface{}
	decodeJSONResponse(t, roomsResponse, &roomsList)
	foundRoom := false
	for _, room := range roomsList {
		if room["id"] == roomID.String() && room["name"] == "Sala Teste" {
			foundRoom = true
			break
		}
	}
	if !foundRoom {
		t.Fatal("expected seeded telemetry room in rooms list")
	}

	bedsResponse := nurseClient.Get(t, "/api/v1/telemetry/rooms/"+roomID.String()+"/beds")
	requireStatusCode(t, bedsResponse, http.StatusOK)
	var bedsList []map[string]interface{}
	decodeJSONResponse(t, bedsResponse, &bedsList)
	if len(bedsList) == 0 {
		t.Fatal("expected seeded telemetry bed in beds list")
	}

	conditionResponse := nurseClient.Post(t, "/api/v1/telemetry/beds/"+bedID.String()+"/condition", map[string]interface{}{
		"bpm":         88,
		"spo2":        96,
		"temperature": 36.8,
		"status":      "normal",
		"condition":   "Normal",
	})
	requireStatusCode(t, conditionResponse, http.StatusOK)

	wrongPasscodeResponse := nurseClient.Post(t, "/api/v1/telemetry/rooms/"+roomID.String()+"/unlock", map[string]string{
		"passcode": "0000",
	})
	requireStatusCode(t, wrongPasscodeResponse, http.StatusForbidden)

	correctPasscodeResponse := nurseClient.Post(t, "/api/v1/telemetry/rooms/"+roomID.String()+"/unlock", map[string]string{
		"passcode": "1234",
	})
	requireStatusCode(t, correctPasscodeResponse, http.StatusOK)
}

func TestPatientCannotAccessTelemetry(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	patientClient := loginAs(t, serverURL, "paciente@mail.com", "secret123")
	requireStatusCode(t, patientClient.Get(t, "/api/v1/telemetry/rooms"), http.StatusForbidden)
}

func seedTelemetryRoomAndBed(t *testing.T, db *pgxpool.Pool) (uuid.UUID, uuid.UUID) {
	t.Helper()
	ctx := context.Background()
	roomID := uuid.New()
	bedID := uuid.New()
	hashedPasscode, hashError := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.MinCost)
	if hashError != nil {
		t.Fatalf("failed to hash telemetry passcode: %v", hashError)
	}
	if _, err := db.Exec(ctx, `
		INSERT INTO telemetry_rooms (id, name, passcode, description)
		VALUES ($1, 'Sala Teste', $2, 'Quarto de teste')`, roomID, string(hashedPasscode)); err != nil {
		t.Fatalf("failed to seed telemetry room: %v", err)
	}
	if _, err := db.Exec(ctx, `
		INSERT INTO telemetry_beds (id, room_id, bed_number, patient_name, age, gender, bpm, spo2, temperature, status, condition)
		VALUES ($1, $2, 'Leito Teste', 'Paciente Teste', 30, 'Feminino', 80, 98, 36.5, 'normal', 'Normal')`, bedID, roomID); err != nil {
		t.Fatalf("failed to seed telemetry bed: %v", err)
	}
	return roomID, bedID
}
