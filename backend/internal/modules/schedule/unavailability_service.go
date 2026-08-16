package schedule

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

type UnavailabilityService interface {
	CreateUnavailability(ctx context.Context, input CreateUnavailabilityInput) (*StaffUnavailability, error)
	ListUnavailabilityByStaff(ctx context.Context, input ListUnavailabilityInput) ([]*StaffUnavailability, error)
	DeleteUnavailability(ctx context.Context, unavailabilityID uuid.UUID) (*StaffUnavailability, error)
}

type unavailabilityService struct {
	repo UnavailabilityRepository
}

func NewUnavailabilityService(repo UnavailabilityRepository) UnavailabilityService {
	return &unavailabilityService{repo: repo}
}

func (unavailabilityService *unavailabilityService) CreateUnavailability(ctx context.Context, input CreateUnavailabilityInput) (*StaffUnavailability, error) {
	fieldViolations := make(map[string]string)
	if input.StaffID == uuid.Nil {
		fieldViolations["staff_id"] = "is required"
	}
	if input.StartsAt.IsZero() {
		fieldViolations["starts_at"] = "is required"
	} else if !input.StartsAt.After(time.Now()) {
		fieldViolations["starts_at"] = "must be in the future"
	}
	if input.EndsAt.IsZero() {
		fieldViolations["ends_at"] = "is required"
	} else if !input.EndsAt.After(input.StartsAt) {
		fieldViolations["ends_at"] = "must be after starts_at"
	}
	if len(input.Reason) > 500 {
		fieldViolations["reason"] = "must be at most 500 characters"
	}
	if len(fieldViolations) > 0 {
		return nil, apperrors.InvalidArgument("invalid unavailability input", fieldViolations)
	}

	newUnavailability := &StaffUnavailability{
		ID:        uuid.New(),
		StaffID:   input.StaffID,
		StartsAt:  input.StartsAt,
		EndsAt:    input.EndsAt,
		Reason:    input.Reason,
		CreatedBy: actorFromContext(ctx),
	}

	return unavailabilityService.repo.CreateUnavailability(ctx, newUnavailability)
}

func (unavailabilityService *unavailabilityService) ListUnavailabilityByStaff(ctx context.Context, input ListUnavailabilityInput) ([]*StaffUnavailability, error) {
	if input.StaffID == uuid.Nil {
		return nil, apperrors.InvalidArgument("invalid unavailability filter", map[string]string{"staff_id": "is required"})
	}
	if !input.From.IsZero() && !input.To.IsZero() && !input.To.After(input.From) {
		return nil, apperrors.InvalidArgument("invalid unavailability filter", map[string]string{"to": "must be after from"})
	}
	return unavailabilityService.repo.ListUnavailabilityByStaff(ctx, input.StaffID, input.From, input.To)
}

func (unavailabilityService *unavailabilityService) DeleteUnavailability(ctx context.Context, unavailabilityID uuid.UUID) (*StaffUnavailability, error) {
	if unavailabilityID == uuid.Nil {
		return nil, apperrors.InvalidArgument("invalid unavailability input", map[string]string{"id": "is required"})
	}
	return unavailabilityService.repo.DeleteUnavailability(ctx, unavailabilityID)
}
