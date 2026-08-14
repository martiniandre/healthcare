package fhir

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestNewConditionResource_OmitsEncounterWhenEmpty(t *testing.T) {
	resource := NewConditionResource("pat-1", "", "J45.9", "Asthma", "active", time.Now())

	encodedBody, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("failed to marshal condition resource: %v", err)
	}
	if strings.Contains(string(encodedBody), `"encounter"`) {
		t.Fatalf("expected encounter to be omitted, got %s", encodedBody)
	}
}

func TestNewConditionResource_IncludesEncounterWhenProvided(t *testing.T) {
	resource := NewConditionResource("pat-1", "enc-9", "J45.9", "Asthma", "active", time.Now())

	encodedBody, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("failed to marshal condition resource: %v", err)
	}
	if !strings.Contains(string(encodedBody), `"reference":"Encounter/enc-9"`) {
		t.Fatalf("expected encounter reference to be present, got %s", encodedBody)
	}
}

func TestNewMedicationRequestResource_OmitsEncounterWhenEmpty(t *testing.T) {
	resource := NewMedicationRequestResource("pat-1", "", "prac-1", "123456", "Amoxicilina", "1 comprimido 8/8h")

	encodedBody, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("failed to marshal medication request resource: %v", err)
	}
	if strings.Contains(string(encodedBody), `"encounter"`) {
		t.Fatalf("expected encounter to be omitted, got %s", encodedBody)
	}
}

func TestNewDiagnosticReportResource_OmitsEncounterWhenEmpty(t *testing.T) {
	resource := NewDiagnosticReportResource("pat-1", "", "11529-2", "Radiografia de tórax", "Sem alterações")

	encodedBody, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("failed to marshal diagnostic report resource: %v", err)
	}
	if strings.Contains(string(encodedBody), `"encounter"`) {
		t.Fatalf("expected encounter to be omitted, got %s", encodedBody)
	}
}
