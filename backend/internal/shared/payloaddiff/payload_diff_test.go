package payloaddiff

import (
	"reflect"
	"testing"
)

func TestCompute_DetectsChangedFieldWithBeforeAndAfter(t *testing.T) {
	before := map[string]any{"conclusion": "old", "status": "preliminary"}
	after := map[string]any{"conclusion": "new", "status": "preliminary"}

	changes, computeErr := Compute(before, after)
	if computeErr != nil {
		t.Fatalf("unexpected error: %v", computeErr)
	}
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	conclusionChange, exists := changes["conclusion"]
	if !exists {
		t.Fatal("expected conclusion key in changes")
	}
	expectedChange := map[string]any{"before": "old", "after": "new"}
	if !reflect.DeepEqual(conclusionChange, expectedChange) {
		t.Errorf("expected %v, got %v", expectedChange, conclusionChange)
	}
}

func TestCompute_ReturnsEmptyMapWhenUnchanged(t *testing.T) {
	before := map[string]any{"conclusion": "same", "status": "final"}
	after := map[string]any{"status": "final", "conclusion": "same"}

	changes, computeErr := Compute(before, after)
	if computeErr != nil {
		t.Fatalf("unexpected error: %v", computeErr)
	}
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %v", changes)
	}
}

func TestCompute_FlagsAddedAndRemovedFields(t *testing.T) {
	before := map[string]any{"existing": "value"}
	after := map[string]any{"new_field": "added"}

	changes, computeErr := Compute(before, after)
	if computeErr != nil {
		t.Fatalf("unexpected error: %v", computeErr)
	}
	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(changes))
	}
	if _, exists := changes["new_field"]; !exists {
		t.Error("expected new_field change")
	}
	if _, exists := changes["existing"]; !exists {
		t.Error("expected existing change")
	}
}

func TestCompute_IsDeterministicAcrossCalls(t *testing.T) {
	before := map[string]any{"b": 2, "a": 1, "nested": map[string]any{"z": "x", "y": "w"}}
	after := map[string]any{"b": 3, "a": 1, "nested": map[string]any{"z": "changed", "y": "w"}}

	firstChanges, firstErr := Compute(before, after)
	secondChanges, secondErr := Compute(before, after)
	if firstErr != nil || secondErr != nil {
		t.Fatalf("unexpected error: %v / %v", firstErr, secondErr)
	}
	if !reflect.DeepEqual(firstChanges, secondChanges) {
		t.Errorf("expected deterministic diff, got %v and %v", firstChanges, secondChanges)
	}
}
