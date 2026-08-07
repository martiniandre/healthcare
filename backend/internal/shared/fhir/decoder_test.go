package fhir

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

const encounterBundleFixture = `{
	"resourceType": "Bundle",
	"entry": [
		{
			"resource": {
				"resourceType": "Encounter",
				"id": "enc-123",
				"status": "finished",
				"subject": {"reference": "Patient/pat-456"},
				"period": {"start": "2026-07-01T10:00:00Z", "end": "2026-07-01T11:30:00Z"},
				"reasonCode": [
					{"coding": [{"system": "http://hl7.org/fhir/sid/icd-10", "code": "R10.9", "display": "Abdominal pain"}], "text": "Unspecified abdominal pain"}
				],
				"participant": [
					{"type": [{"coding": [{"code": "PPRF"}]}], "individual": {"reference": "Practitioner/doc-789"}}
				]
			}
		},
		{
			"resource": {
				"resourceType": "Encounter",
				"id": "enc-124",
				"status": "in-progress",
				"subject": {"reference": "Patient/pat-456"},
				"period": {"start": "2026-07-02T09:00:00Z"},
				"reasonCode": [
					{"coding": [{"code": "A00", "display": "Cholera"}]}
				]
			}
		}
	]
}`

const observationBundleFixture = `{
	"resourceType": "Bundle",
	"entry": [
		{
			"resource": {
				"resourceType": "Observation",
				"id": "obs-1",
				"status": "final",
				"code": {
					"coding": [{"system": "http://loinc.org", "code": "718-7", "display": "Hemoglobin [Mass/volume] in Blood"}],
					"text": "Hemoglobin"
				},
				"subject": {"reference": "Patient/pat-456"},
				"encounter": {"reference": "Encounter/enc-123"},
				"valueQuantity": {"value": 14.2, "unit": "g/dL"},
				"effectiveDateTime": "2026-07-01T10:05:00Z",
				"issued": "2026-07-01T10:15:00Z"
			}
		}
	]
}`

const conditionBundleFixture = `{
	"resourceType": "Bundle",
	"entry": [
		{
			"resource": {
				"resourceType": "Condition",
				"id": "cond-1",
				"clinicalStatus": {
					"coding": [{"system": "http://terminology.hl7.org/CodeSystem/condition-clinical", "code": "active", "display": "Active"}]
				},
				"code": {
					"coding": [{"system": "http://hl7.org/fhir/sid/icd-10", "code": "E11.9", "display": "Type 2 diabetes mellitus"}],
					"text": "Diabetes"
				},
				"subject": {"reference": "Patient/pat-456"},
				"encounter": {"reference": "Encounter/enc-123"},
				"onsetDateTime": "2026-05-01T00:00:00Z",
				"recordedDate": "2026-07-01T10:10:00Z"
			}
		}
	]
}`

const allergyBundleFixture = `{
	"resourceType": "Bundle",
	"entry": [
		{
			"resource": {
				"resourceType": "AllergyIntolerance",
				"id": "all-1",
				"clinicalStatus": {
					"coding": [{"system": "http://terminology.hl7.org/CodeSystem/allergyintolerance-clinical", "code": "active", "display": "Active"}]
				},
				"code": {
					"coding": [{"system": "http://www.nlm.nih.gov/research/umls/rxnorm", "code": "7980", "display": "Penicillin"}],
					"text": "Penicillin allergy"
				},
				"patient": {"reference": "Patient/pat-456"},
				"recordedDate": "2026-06-20T00:00:00Z",
				"reaction": [
					{"manifestation": [{"text": "Urticaria", "coding": [{"code": "247472004"}]}]}
				]
			}
		}
	]
}`

const medicationRequestBundleFixture = `{
	"resourceType": "Bundle",
	"entry": [
		{
			"resource": {
				"resourceType": "MedicationRequest",
				"id": "med-1",
				"status": "active",
				"intent": "order",
				"medicationCodeableConcept": {
					"coding": [{"system": "http://www.nlm.nih.gov/research/umls/rxnorm", "code": "6013", "display": "Metformin 500 MG"}],
					"text": "Metformin"
				},
				"subject": {"reference": "Patient/pat-456"},
				"encounter": {"reference": "Encounter/enc-123"},
				"requester": {"agent": {"reference": "Practitioner/doc-789"}},
				"authoredOn": "2026-07-01T10:20:00Z",
				"dosageInstruction": [{"text": "Take one tablet twice daily with meals"}]
			}
		}
	]
}`

const diagnosticReportBundleFixture = `{
	"resourceType": "Bundle",
	"entry": [
		{
			"resource": {
				"resourceType": "DiagnosticReport",
				"id": "rep-1",
				"status": "final",
				"meta": {"versionId": "3", "lastUpdated": "2026-07-01T11:00:00Z"},
				"code": {
					"coding": [{"system": "http://loinc.org", "code": "58410-2", "display": "Complete blood count"}],
					"text": "Blood panel"
				},
				"subject": {"reference": "Patient/pat-456"},
				"encounter": {"reference": "Encounter/enc-123"},
				"issued": "2026-07-01T11:05:00Z",
				"conclusion": "Within normal limits"
			}
		}
	]
}`

const patientResourceFixture = `{
	"resourceType": "Patient",
	"id": "pat-456",
	"meta": {"versionId": "2", "lastUpdated": "2026-07-01T12:00:00Z"},
	"name": [{"use": "official", "family": "Silva", "given": ["Maria"]}],
	"birthDate": "1985-03-14",
	"telecom": [{"system": "phone", "value": "+5511999999999", "use": "mobile"}],
	"identifier": [{"system": "urn:oid:2.16.840.1.113883.13.237", "value": "12345678900"}]
}`

const imagingStudyBundleFixture = `{
	"resourceType": "Bundle",
	"entry": [
		{
			"resource": {
				"resourceType": "ImagingStudy",
				"id": "img-1",
				"status": "available",
				"started": "2026-07-01T13:00:00Z",
				"description": "Chest X-ray",
				"modality": [{"system": "http://dicom.nema.org/resources/ontology/DCM", "code": "CR", "display": "Computed Radiography"}]
			}
		}
	]
}`

func TestDecodeBundleReturnsResourceEntries(t *testing.T) {
	encounters, err := DecodeBundle[Encounter](json.RawMessage(encounterBundleFixture))
	if err != nil {
		t.Fatalf("DecodeBundle returned error: %v", err)
	}
	if len(encounters) != 2 {
		t.Fatalf("expected 2 encounters, got %d", len(encounters))
	}
	first := encounters[0]
	if first.ID != "enc-123" {
		t.Errorf("expected id enc-123, got %q", first.ID)
	}
	if first.Status != "finished" {
		t.Errorf("expected status finished, got %q", first.Status)
	}
	if first.Subject.Reference != "Patient/pat-456" {
		t.Errorf("expected subject Patient/pat-456, got %q", first.Subject.Reference)
	}
	if first.Period == nil {
		t.Fatal("expected period to be decoded")
	}
	if first.Period.Start != "2026-07-01T10:00:00Z" {
		t.Errorf("expected period start, got %q", first.Period.Start)
	}
	if first.Period.End != "2026-07-01T11:30:00Z" {
		t.Errorf("expected period end, got %q", first.Period.End)
	}
	if len(first.ReasonCode) != 1 {
		t.Fatalf("expected 1 reason code, got %d", len(first.ReasonCode))
	}
	if first.ReasonCode[0].Text != "Unspecified abdominal pain" {
		t.Errorf("expected reason text, got %q", first.ReasonCode[0].Text)
	}
	if len(first.Participant) != 1 {
		t.Fatalf("expected 1 participant, got %d", len(first.Participant))
	}
	if first.Participant[0].Individual.Reference != "Practitioner/doc-789" {
		t.Errorf("expected participant reference, got %q", first.Participant[0].Individual.Reference)
	}
	second := encounters[1]
	if second.Period == nil {
		t.Fatal("expected second encounter period to be decoded")
	}
	if second.Period.End != "" {
		t.Errorf("expected empty period end, got %q", second.Period.End)
	}
}

func TestDecodeBundleOnEmptyEntryList(t *testing.T) {
	bundle := []byte(`{"resourceType": "Bundle", "entry": []}`)
	resources, err := DecodeBundle[Encounter](bundle)
	if err != nil {
		t.Fatalf("DecodeBundle returned error: %v", err)
	}
	if len(resources) != 0 {
		t.Errorf("expected 0 resources, got %d", len(resources))
	}
}

func TestDecodeBundleOnMalformedBundle(t *testing.T) {
	_, err := DecodeBundle[Encounter]([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for malformed bundle")
	}
}

func TestDecodeBundleOnEmptyResourceEntry(t *testing.T) {
	bundle := []byte(`{"resourceType": "Bundle", "entry": [{"resource": null}]}`)
	resources, err := DecodeBundle[Encounter](bundle)
	if err != nil {
		t.Fatalf("DecodeBundle returned error: %v", err)
	}
	if len(resources) != 0 {
		t.Errorf("expected 0 resources for null entries, got %d", len(resources))
	}
}

func TestDecodeResourceDecodesObservation(t *testing.T) {
	resource, err := DecodeResource[Observation](json.RawMessage(`{
		"resourceType": "Observation",
		"id": "obs-1",
		"code": {"coding": [{"code": "718-7", "display": "Hemoglobin"}], "text": "Hemoglobin"},
		"subject": {"reference": "Patient/pat-456"},
		"encounter": {"reference": "Encounter/enc-123"},
		"valueQuantity": {"value": 14.2, "unit": "g/dL"},
		"effectiveDateTime": "2026-07-01T10:05:00Z"
	}`))
	if err != nil {
		t.Fatalf("DecodeResource returned error: %v", err)
	}
	if resource.ID != "obs-1" {
		t.Errorf("expected id obs-1, got %q", resource.ID)
	}
	if resource.ValueQuantity == nil {
		t.Fatal("expected valueQuantity to be decoded")
	}
	if resource.ValueQuantity.Value != 14.2 {
		t.Errorf("expected value 14.2, got %v", resource.ValueQuantity.Value)
	}
	if resource.ValueQuantity.Unit != "g/dL" {
		t.Errorf("expected unit g/dL, got %q", resource.ValueQuantity.Unit)
	}
}

func TestDecodeResourceOnMalformedResource(t *testing.T) {
	_, err := DecodeResource[Patient]([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for malformed resource")
	}
}

func TestCodeableConceptPartsExtractsCodeDisplayAndText(t *testing.T) {
	concept := CodeableConcept{
		Coding: []Coding{{Code: "718-7", Display: "Hemoglobin"}},
		Text:   "Hemoglobin panel",
	}
	code, display, text := CodeableConceptParts(concept)
	if code != "718-7" {
		t.Errorf("expected code 718-7, got %q", code)
	}
	if display != "Hemoglobin" {
		t.Errorf("expected display Hemoglobin, got %q", display)
	}
	if text != "Hemoglobin panel" {
		t.Errorf("expected text Hemoglobin panel, got %q", text)
	}
}

func TestCodeableConceptPartsOnEmptyConcept(t *testing.T) {
	code, display, text := CodeableConceptParts(CodeableConcept{})
	if code != "" || display != "" || text != "" {
		t.Errorf("expected empty parts, got code=%q display=%q text=%q", code, display, text)
	}
}

func TestSplitReferenceIDSplitsResourceTypePrefix(t *testing.T) {
	if got := SplitReferenceID("Patient/pat-456"); got != "pat-456" {
		t.Errorf("expected pat-456, got %q", got)
	}
	if got := SplitReferenceID("Encounter/enc-123"); got != "enc-123" {
		t.Errorf("expected enc-123, got %q", got)
	}
}

func TestSplitReferenceIDWithoutSeparator(t *testing.T) {
	if got := SplitReferenceID("pat-456"); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestParseRFC3339ParsesValidTimestamp(t *testing.T) {
	parsed, ok := ParseRFC3339("2026-07-01T10:00:00Z")
	if !ok {
		t.Fatal("expected valid parse")
	}
	expected := time.Date(2026, 7, 1, 10, 0, 0, 0, time.UTC)
	if !parsed.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, parsed)
	}
}

func TestParseRFC3339RejectsInvalidTimestamp(t *testing.T) {
	if _, ok := ParseRFC3339("not-a-date"); ok {
		t.Fatal("expected invalid parse")
	}
}

func TestResourceVersionIDFromMeta(t *testing.T) {
	if got := ResourceVersionID(&ResourceMeta{VersionID: "3"}); got != "3" {
		t.Errorf("expected 3, got %q", got)
	}
	if got := ResourceVersionID(&ResourceMeta{}); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
	if got := ResourceVersionID(nil); got != "" {
		t.Errorf("expected empty for nil meta, got %q", got)
	}
}

func TestDecodeBundleWithMissingKeysLeavesZeroValues(t *testing.T) {
	bundle := []byte(`{"resourceType": "Bundle", "entry": [{"resource": {"id": "x"}}]}`)
	resources, err := DecodeBundle[Encounter](bundle)
	if err != nil {
		t.Fatalf("DecodeBundle returned error: %v", err)
	}
	if len(resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(resources))
	}
	if resources[0].Status != "" {
		t.Errorf("expected empty status, got %q", resources[0].Status)
	}
	if resources[0].Subject.Reference != "" {
		t.Errorf("expected empty subject reference, got %q", resources[0].Subject.Reference)
	}
}

func TestFieldLevelEqualityForObservationParserShape(t *testing.T) {
	observations, err := DecodeBundle[Observation](json.RawMessage(observationBundleFixture))
	if err != nil {
		t.Fatalf("DecodeBundle returned error: %v", err)
	}
	if len(observations) != 1 {
		t.Fatalf("expected 1 observation, got %d", len(observations))
	}
	observation := observations[0]
	if observation.ID != "obs-1" {
		t.Errorf("expected id obs-1, got %q", observation.ID)
	}
	code, display, text := CodeableConceptParts(observation.Code)
	if code != "718-7" {
		t.Errorf("expected loinc 718-7, got %q", code)
	}
	if display != "Hemoglobin [Mass/volume] in Blood" {
		t.Errorf("expected display, got %q", display)
	}
	if text != "Hemoglobin" {
		t.Errorf("expected text Hemoglobin, got %q", text)
	}
	if observation.Encounter.Reference != "Encounter/enc-123" {
		t.Errorf("expected encounter reference, got %q", observation.Encounter.Reference)
	}
	if observation.Subject.Reference != "Patient/pat-456" {
		t.Errorf("expected subject reference, got %q", observation.Subject.Reference)
	}
}

func TestFieldLevelEqualityForConditionParserShape(t *testing.T) {
	conditions, err := DecodeBundle[Condition](json.RawMessage(conditionBundleFixture))
	if err != nil {
		t.Fatalf("DecodeBundle returned error: %v", err)
	}
	if len(conditions) != 1 {
		t.Fatalf("expected 1 condition, got %d", len(conditions))
	}
	condition := conditions[0]
	if condition.ID != "cond-1" {
		t.Errorf("expected id cond-1, got %q", condition.ID)
	}
	clinicalStatusCode, _, _ := CodeableConceptParts(condition.ClinicalStatus)
	if clinicalStatusCode != "active" {
		t.Errorf("expected clinical status active, got %q", clinicalStatusCode)
	}
	code, _, text := CodeableConceptParts(condition.Code)
	if code != "E11.9" {
		t.Errorf("expected icd10 E11.9, got %q", code)
	}
	if text != "Diabetes" {
		t.Errorf("expected text Diabetes, got %q", text)
	}
	if condition.OnsetDateTime != "2026-05-01T00:00:00Z" {
		t.Errorf("expected onset datetime, got %q", condition.OnsetDateTime)
	}
	if condition.RecordedDate != "2026-07-01T10:10:00Z" {
		t.Errorf("expected recorded date, got %q", condition.RecordedDate)
	}
}

func TestFieldLevelEqualityForAllergyParserShape(t *testing.T) {
	allergies, err := DecodeBundle[AllergyIntolerance](json.RawMessage(allergyBundleFixture))
	if err != nil {
		t.Fatalf("DecodeBundle returned error: %v", err)
	}
	if len(allergies) != 1 {
		t.Fatalf("expected 1 allergy, got %d", len(allergies))
	}
	allergy := allergies[0]
	if allergy.ID != "all-1" {
		t.Errorf("expected id all-1, got %q", allergy.ID)
	}
	clinicalStatusCode, _, _ := CodeableConceptParts(allergy.ClinicalStatus)
	if clinicalStatusCode != "active" {
		t.Errorf("expected clinical status active, got %q", clinicalStatusCode)
	}
	code, display, text := CodeableConceptParts(allergy.Code)
	if code != "7980" {
		t.Errorf("expected rxnorm 7980, got %q", code)
	}
	if display != "Penicillin" {
		t.Errorf("expected display Penicillin, got %q", display)
	}
	if text != "Penicillin allergy" {
		t.Errorf("expected text, got %q", text)
	}
	if allergy.Patient.Reference != "Patient/pat-456" {
		t.Errorf("expected patient reference, got %q", allergy.Patient.Reference)
	}
	if len(allergy.Reaction) != 1 || len(allergy.Reaction[0].Manifestation) != 1 {
		t.Fatalf("expected 1 reaction with 1 manifestation")
	}
	if allergy.Reaction[0].Manifestation[0].Text != "Urticaria" {
		t.Errorf("expected manifestation text Urticaria, got %q", allergy.Reaction[0].Manifestation[0].Text)
	}
}

func TestFieldLevelEqualityForMedicationRequestParserShape(t *testing.T) {
	medications, err := DecodeBundle[MedicationRequest](json.RawMessage(medicationRequestBundleFixture))
	if err != nil {
		t.Fatalf("DecodeBundle returned error: %v", err)
	}
	if len(medications) != 1 {
		t.Fatalf("expected 1 medication, got %d", len(medications))
	}
	medication := medications[0]
	if medication.ID != "med-1" {
		t.Errorf("expected id med-1, got %q", medication.ID)
	}
	if medication.Status != "active" {
		t.Errorf("expected status active, got %q", medication.Status)
	}
	code, display, text := CodeableConceptParts(medication.MedicationCodeableConcept)
	if code != "6013" {
		t.Errorf("expected rxnorm 6013, got %q", code)
	}
	if display != "Metformin 500 MG" {
		t.Errorf("expected display, got %q", display)
	}
	if text != "Metformin" {
		t.Errorf("expected text Metformin, got %q", text)
	}
	if medication.Subject.Reference != "Patient/pat-456" {
		t.Errorf("expected subject reference, got %q", medication.Subject.Reference)
	}
	if medication.Requester.Agent.Reference != "Practitioner/doc-789" {
		t.Errorf("expected requester agent reference, got %q", medication.Requester.Agent.Reference)
	}
	if medication.AuthoredOn != "2026-07-01T10:20:00Z" {
		t.Errorf("expected authoredOn, got %q", medication.AuthoredOn)
	}
	if len(medication.DosageInstruction) != 1 {
		t.Fatalf("expected 1 dosage instruction, got %d", len(medication.DosageInstruction))
	}
	if medication.DosageInstruction[0].Text != "Take one tablet twice daily with meals" {
		t.Errorf("expected dosage text, got %q", medication.DosageInstruction[0].Text)
	}
}

func TestFieldLevelEqualityForDiagnosticReportParserShape(t *testing.T) {
	reports, err := DecodeBundle[DiagnosticReport](json.RawMessage(diagnosticReportBundleFixture))
	if err != nil {
		t.Fatalf("DecodeBundle returned error: %v", err)
	}
	if len(reports) != 1 {
		t.Fatalf("expected 1 report, got %d", len(reports))
	}
	report := reports[0]
	if report.ID != "rep-1" {
		t.Errorf("expected id rep-1, got %q", report.ID)
	}
	if report.Status != "final" {
		t.Errorf("expected status final, got %q", report.Status)
	}
	if report.Meta == nil {
		t.Fatal("expected meta to be decoded")
	}
	if report.Meta.VersionID != "3" {
		t.Errorf("expected versionId 3, got %q", report.Meta.VersionID)
	}
	if report.Meta.LastUpdated != "2026-07-01T11:00:00Z" {
		t.Errorf("expected lastUpdated, got %q", report.Meta.LastUpdated)
	}
	code, display, text := CodeableConceptParts(report.Code)
	if code != "58410-2" {
		t.Errorf("expected loinc 58410-2, got %q", code)
	}
	if display != "Complete blood count" {
		t.Errorf("expected display, got %q", display)
	}
	if text != "Blood panel" {
		t.Errorf("expected text, got %q", text)
	}
	if report.Issued != "2026-07-01T11:05:00Z" {
		t.Errorf("expected issued, got %q", report.Issued)
	}
	if report.Conclusion != "Within normal limits" {
		t.Errorf("expected conclusion, got %q", report.Conclusion)
	}
}

func TestFieldLevelEqualityForPatientParserShape(t *testing.T) {
	patient, err := DecodeResource[Patient](json.RawMessage(patientResourceFixture))
	if err != nil {
		t.Fatalf("DecodeResource returned error: %v", err)
	}
	if patient.ID != "pat-456" {
		t.Errorf("expected id pat-456, got %q", patient.ID)
	}
	if patient.Meta == nil {
		t.Fatal("expected meta to be decoded")
	}
	if patient.Meta.LastUpdated != "2026-07-01T12:00:00Z" {
		t.Errorf("expected lastUpdated, got %q", patient.Meta.LastUpdated)
	}
	if len(patient.Name) != 1 {
		t.Fatalf("expected 1 name, got %d", len(patient.Name))
	}
	if patient.Name[0].Family != "Silva" {
		t.Errorf("expected family Silva, got %q", patient.Name[0].Family)
	}
	if len(patient.Name[0].Given) != 1 || patient.Name[0].Given[0] != "Maria" {
		t.Errorf("expected given Maria, got %v", patient.Name[0].Given)
	}
	if patient.BirthDate != "1985-03-14" {
		t.Errorf("expected birthDate, got %q", patient.BirthDate)
	}
	if len(patient.Telecom) != 1 {
		t.Fatalf("expected 1 telecom, got %d", len(patient.Telecom))
	}
	if patient.Telecom[0].System != "phone" || patient.Telecom[0].Value != "+5511999999999" {
		t.Errorf("expected phone telecom, got %+v", patient.Telecom[0])
	}
	if len(patient.Identifier) != 1 {
		t.Fatalf("expected 1 identifier, got %d", len(patient.Identifier))
	}
	if patient.Identifier[0].Value != "12345678900" {
		t.Errorf("expected identifier value, got %q", patient.Identifier[0].Value)
	}
}

func TestFieldLevelEqualityForImagingStudyParserShape(t *testing.T) {
	imagingStudies, err := DecodeBundle[ImagingStudy](json.RawMessage(imagingStudyBundleFixture))
	if err != nil {
		t.Fatalf("DecodeBundle returned error: %v", err)
	}
	if len(imagingStudies) != 1 {
		t.Fatalf("expected 1 imaging study, got %d", len(imagingStudies))
	}
	imagingStudy := imagingStudies[0]
	if imagingStudy.ID != "img-1" {
		t.Errorf("expected id img-1, got %q", imagingStudy.ID)
	}
	if imagingStudy.Status != "available" {
		t.Errorf("expected status available, got %q", imagingStudy.Status)
	}
	if imagingStudy.Started != "2026-07-01T13:00:00Z" {
		t.Errorf("expected started, got %q", imagingStudy.Started)
	}
	if imagingStudy.Description != "Chest X-ray" {
		t.Errorf("expected description, got %q", imagingStudy.Description)
	}
	if len(imagingStudy.Modality) != 1 {
		t.Fatalf("expected 1 modality, got %d", len(imagingStudy.Modality))
	}
	if imagingStudy.Modality[0].Code != "CR" {
		t.Errorf("expected modality code CR, got %q", imagingStudy.Modality[0].Code)
	}
	if imagingStudy.Modality[0].Display != "Computed Radiography" {
		t.Errorf("expected modality display, got %q", imagingStudy.Modality[0].Display)
	}
}

func TestDecodeBundleResilienceToMissingResourceType(t *testing.T) {
	bundle := []byte(`{"resourceType": "Bundle", "entry": [{"resource": {"id": "no-type"}}]}`)
	resources, err := DecodeBundle[Condition](bundle)
	if err != nil {
		t.Fatalf("DecodeBundle returned error: %v", err)
	}
	if len(resources) != 1 || resources[0].ID != "no-type" {
		t.Errorf("expected resource decoded without resourceType, got %+v", resources)
	}
	if resources[0].ResourceType != "" {
		t.Errorf("expected empty resource type, got %q", resources[0].ResourceType)
	}
}

func TestConcatenatedNameFormatting(t *testing.T) {
	family := "Silva"
	given := "Maria"
	trimmed := strings.TrimSpace(given + " " + family)
	if trimmed != "Maria Silva" {
		t.Errorf("expected Maria Silva, got %q", trimmed)
	}
}
