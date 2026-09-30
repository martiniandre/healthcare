package telemetry

import "strings"

type BedStatus string

const (
	BedStatusNormal  BedStatus = "normal"
	BedStatusWarning BedStatus = "warning"
	BedStatusDanger  BedStatus = "danger"
)

type CardiacCondition string

const (
	CardiacConditionNormal        CardiacCondition = "normal"
	CardiacConditionBradycardia   CardiacCondition = "bradycardia"
	CardiacConditionTachycardia   CardiacCondition = "tachycardia"
	CardiacConditionCardiacArrest CardiacCondition = "cardiac-arrest"
)

var legacyCardiacConditionAliases = map[string]CardiacCondition{
	"bradicardia":     CardiacConditionBradycardia,
	"taquicardia":     CardiacConditionTachycardia,
	"parada cardiaca": CardiacConditionCardiacArrest,
	"parada cardíaca": CardiacConditionCardiacArrest,
}

func NormalizeBedStatus(rawStatus string) BedStatus {
	switch strings.ToLower(strings.TrimSpace(rawStatus)) {
	case string(BedStatusWarning):
		return BedStatusWarning
	case string(BedStatusDanger):
		return BedStatusDanger
	default:
		return BedStatusNormal
	}
}

func NormalizeCardiacCondition(rawCondition string) CardiacCondition {
	normalizedKey := strings.ToLower(strings.TrimSpace(rawCondition))
	switch normalizedKey {
	case string(CardiacConditionBradycardia):
		return CardiacConditionBradycardia
	case string(CardiacConditionTachycardia):
		return CardiacConditionTachycardia
	case string(CardiacConditionCardiacArrest):
		return CardiacConditionCardiacArrest
	}
	if legacyCondition, aliasFound := legacyCardiacConditionAliases[normalizedKey]; aliasFound {
		return legacyCondition
	}
	return CardiacConditionNormal
}
