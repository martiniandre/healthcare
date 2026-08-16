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

func TestNewEncounterResource_WithCodeAndDisplay_PersistsCodingAndText(t *testing.T) {
	resource := NewEncounterResource("pat-1", "prac-1", "Z00.0", "Routine check-up")

	if len(resource.ReasonCode) != 1 {
		t.Fatalf("expected 1 reason code, got %d", len(resource.ReasonCode))
	}
	reasonConcept := resource.ReasonCode[0]
	if len(reasonConcept.Coding) != 1 {
		t.Fatalf("expected 1 coding, got %d", len(reasonConcept.Coding))
	}
	if reasonConcept.Coding[0].System != "http://hl7.org/fhir/sid/icd-10" {
		t.Errorf("expected icd-10 system, got %q", reasonConcept.Coding[0].System)
	}
	if reasonConcept.Coding[0].Code != "Z00.0" {
		t.Errorf("expected icd code Z00.0, got %q", reasonConcept.Coding[0].Code)
	}
	if reasonConcept.Coding[0].Display != "Routine check-up" {
		t.Errorf("expected display Routine check-up, got %q", reasonConcept.Coding[0].Display)
	}
	if reasonConcept.Text != "Routine check-up" {
		t.Errorf("expected text Routine check-up, got %q", reasonConcept.Text)
	}
}

func TestNewEncounterResource_WithDisplayOnly_PersistsTextWithoutCoding(t *testing.T) {
	resource := NewEncounterResource("pat-1", "prac-1", "", "Dor abdominal persistente")

	if len(resource.ReasonCode) != 1 {
		t.Fatalf("expected 1 reason code, got %d", len(resource.ReasonCode))
	}
	reasonConcept := resource.ReasonCode[0]
	if len(reasonConcept.Coding) != 0 {
		t.Fatalf("expected no coding, got %d", len(reasonConcept.Coding))
	}
	if reasonConcept.Text != "Dor abdominal persistente" {
		t.Errorf("expected text Dor abdominal persistente, got %q", reasonConcept.Text)
	}
}

func TestNewEncounterResource_WithoutReason_OmitsReasonCode(t *testing.T) {
	resource := NewEncounterResource("pat-1", "prac-1", "", "")

	if len(resource.ReasonCode) != 0 {
		t.Fatalf("expected no reason code, got %d", len(resource.ReasonCode))
	}
}

func TestNewObservationResource_OmitsEncounterWhenEmpty(t *testing.T) {
	resource := NewObservationResource("pat-1", "", "8867-4", "Frequência cardíaca", 72, "beats/minute")

	encodedBody, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("failed to marshal observation resource: %v", err)
	}
	if strings.Contains(string(encodedBody), `"encounter"`) {
		t.Fatalf("expected encounter to be omitted, got %s", encodedBody)
	}
}

func TestNewObservationResource_IncludesEncounterWhenProvided(t *testing.T) {
	resource := NewObservationResource("pat-1", "enc-9", "8867-4", "Frequência cardíaca", 72, "beats/minute")

	encodedBody, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("failed to marshal observation resource: %v", err)
	}
	if !strings.Contains(string(encodedBody), `"reference":"Encounter/enc-9"`) {
		t.Fatalf("expected encounter reference to be present, got %s", encodedBody)
	}
}
