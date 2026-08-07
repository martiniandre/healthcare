package fhir

import (
	"encoding/json"
	"strings"
	"time"
)

type bundleEnvelope struct {
	Entry []struct {
		Resource json.RawMessage `json:"resource"`
	} `json:"entry"`
}

func DecodeBundle[T any](responseBody json.RawMessage) ([]T, error) {
	var envelope bundleEnvelope
	if err := json.Unmarshal(responseBody, &envelope); err != nil {
		return nil, err
	}
	resources := make([]T, 0, len(envelope.Entry))
	for _, entry := range envelope.Entry {
		if len(entry.Resource) == 0 || string(entry.Resource) == "null" {
			continue
		}
		var resource T
		if err := json.Unmarshal(entry.Resource, &resource); err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}
	return resources, nil
}

func DecodeResource[T any](responseBody json.RawMessage) (*T, error) {
	var resource T
	if err := json.Unmarshal(responseBody, &resource); err != nil {
		return nil, err
	}
	return &resource, nil
}

func SplitReferenceID(reference string) string {
	parts := strings.SplitN(reference, "/", 2)
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

func CodeableConceptParts(concept CodeableConcept) (code string, display string, text string) {
	if len(concept.Coding) > 0 {
		code = concept.Coding[0].Code
		display = concept.Coding[0].Display
	}
	return code, display, concept.Text
}

func ParseRFC3339(value string) (time.Time, bool) {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, false
	}
	return parsed, true
}

func ResourceVersionID(meta *ResourceMeta) string {
	if meta == nil {
		return ""
	}
	return meta.VersionID
}
