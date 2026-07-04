package encounter

import (
	"context"
	"errors"

	pb "github.com/healthcare/backend/internal/modules/encounter/pb"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

type GRPCHandler struct {
	service Service
}

func NewGRPCHandler(service Service) *GRPCHandler {
	return &GRPCHandler{service: service}
}

func mapEncounterError(err error) error {
	if errors.Is(err, ErrEncounterNotFound) {
		return apperrors.ErrEncounterNotFound.ToGRPC()
	}
	return apperrors.ToGRPCStatus(err)
}

func (handler *GRPCHandler) CreateEncounter(ctx context.Context, req *pb.CreateEncounterRequest) (*pb.CreateEncounterResponse, error) {
	violations := make(map[string]string)
	if req.PatientFhirId == "" {
		violations["patient_fhir_id"] = "is required"
	}
	if req.PractitionerId == "" {
		violations["practitioner_id"] = "is required"
	}
	if len(violations) > 0 {
		return nil, apperrors.ErrBadRequest.WithFields(violations)
	}

	encounter := &Encounter{
		PatientFHIRID:  req.PatientFhirId,
		PractitionerID: req.PractitionerId,
		ReasonCode:     req.ReasonCode,
		ReasonDisplay:  req.ReasonDisplay,
		Status:         "in-progress",
	}

	createdEncounter, err := handler.service.CreateEncounter(ctx, encounter)
	if err != nil {
		return nil, mapEncounterError(err)
	}

	return &pb.CreateEncounterResponse{EncounterFhirId: createdEncounter.FHIRResourceID}, nil
}

func (handler *GRPCHandler) GetEncounter(ctx context.Context, req *pb.GetEncounterRequest) (*pb.GetEncounterResponse, error) {
	if req.EncounterFhirId == "" {
		return nil, apperrors.ErrBadRequest.WithFields(map[string]string{"encounter_fhir_id": "is required"})
	}

	encounter, err := handler.service.GetEncounter(ctx, req.EncounterFhirId)
	if err != nil {
		return nil, mapEncounterError(err)
	}

	return &pb.GetEncounterResponse{
		FhirId:        encounter.FHIRResourceID,
		PatientFhirId: encounter.PatientFHIRID,
		Status:        encounter.Status,
		ReasonDisplay: encounter.ReasonDisplay,
	}, nil
}

func (handler *GRPCHandler) GetEncounters(ctx context.Context, req *pb.GetEncountersRequest) (*pb.GetEncountersResponse, error) {
	if req.PatientFhirId == "" {
		return nil, apperrors.ErrBadRequest.WithFields(map[string]string{"patient_fhir_id": "is required"})
	}

	encounters, err := handler.service.GetEncountersByPatient(ctx, req.PatientFhirId)
	if err != nil {
		return nil, mapEncounterError(err)
	}

	pbEncounters := make([]*pb.Encounter, 0, len(encounters))
	for _, encounter := range encounters {
		pbEncounters = append(pbEncounters, &pb.Encounter{
			FhirId:        encounter.FHIRResourceID,
			PatientFhirId: encounter.PatientFHIRID,
			Status:        encounter.Status,
			ReasonDisplay: encounter.ReasonDisplay,
		})
	}

	return &pb.GetEncountersResponse{Encounters: pbEncounters}, nil
}
