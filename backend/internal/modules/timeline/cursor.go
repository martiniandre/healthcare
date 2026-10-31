package timeline

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

var errInvalidCursor = errors.New("invalid cursor")

type cursorPosition struct {
	OccurredAt time.Time `json:"t"`
	Type       EntryType `json:"y"`
	ResourceID string    `json:"i"`
}

func encodeCursor(entry Entry) string {
	encoded, marshalError := json.Marshal(cursorPosition{
		OccurredAt: entry.OccurredAt.UTC(),
		Type:       entry.Type,
		ResourceID: entry.ResourceID,
	})
	if marshalError != nil {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(encoded)
}

func decodeCursor(rawCursor string) (cursorPosition, error) {
	decoded, decodeError := base64.RawURLEncoding.DecodeString(rawCursor)
	if decodeError != nil {
		return cursorPosition{}, errInvalidCursor
	}
	var position cursorPosition
	if unmarshalError := json.Unmarshal(decoded, &position); unmarshalError != nil {
		return cursorPosition{}, errInvalidCursor
	}
	if position.OccurredAt.IsZero() || position.ResourceID == "" || !isKnownEntryType(position.Type) {
		return cursorPosition{}, errInvalidCursor
	}
	return position, nil
}

func isKnownEntryType(entryType EntryType) bool {
	for _, knownType := range AllEntryTypes {
		if knownType == entryType {
			return true
		}
	}
	return false
}

func entryComesBefore(first Entry, second Entry) bool {
	if !first.OccurredAt.Equal(second.OccurredAt) {
		return first.OccurredAt.After(second.OccurredAt)
	}
	if first.Type != second.Type {
		return first.Type < second.Type
	}
	return first.ResourceID < second.ResourceID
}

func entryComesAfterCursor(entry Entry, position cursorPosition) bool {
	cursorEntry := Entry{Type: position.Type, ResourceID: position.ResourceID, OccurredAt: position.OccurredAt}
	return entryComesBefore(cursorEntry, entry)
}
