package medication

import "github.com/healthcare/backend/internal/shared/apperrors"

type BlockedByAllergyError struct {
	Alerts []SafetyAlert
}

func (blockedError *BlockedByAllergyError) Error() string {
	return apperrors.ErrMedicationBlockedByAllergy.Message
}

func (blockedError *BlockedByAllergyError) Unwrap() error {
	return apperrors.ErrMedicationBlockedByAllergy
}
