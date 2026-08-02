package payloaddiff

import (
	"encoding/json"
	"reflect"
)

func Compute(before any, after any) (map[string]any, error) {
	beforeMap, beforeErr := Normalize(before)
	if beforeErr != nil {
		return nil, beforeErr
	}
	afterMap, afterErr := Normalize(after)
	if afterErr != nil {
		return nil, afterErr
	}

	changes := make(map[string]any)
	allKeys := make(map[string]struct{}, len(beforeMap)+len(afterMap))
	for key := range beforeMap {
		allKeys[key] = struct{}{}
	}
	for key := range afterMap {
		allKeys[key] = struct{}{}
	}

	orderedKeys := make([]string, 0, len(allKeys))
	for key := range allKeys {
		orderedKeys = append(orderedKeys, key)
	}

	for _, key := range orderedKeys {
		beforeValue, beforeExists := beforeMap[key]
		afterValue, afterExists := afterMap[key]
		if beforeExists == afterExists && beforeExists && reflect.DeepEqual(beforeValue, afterValue) {
			continue
		}
		change := make(map[string]any)
		if beforeExists {
			change["before"] = beforeValue
		}
		if afterExists {
			change["after"] = afterValue
		}
		changes[key] = change
	}

	return changes, nil
}

func Normalize(value any) (map[string]any, error) {
	if value == nil {
		return map[string]any{}, nil
	}
	marshaledValue, marshalErr := json.Marshal(value)
	if marshalErr != nil {
		return nil, marshalErr
	}
	normalizedMap := make(map[string]any)
	if unmarshalErr := json.Unmarshal(marshaledValue, &normalizedMap); unmarshalErr != nil {
		return nil, unmarshalErr
	}
	return normalizedMap, nil
}
