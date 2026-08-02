package encounter

import (
	"context"

	pb "github.com/healthcare/backend/internal/modules/encounter/pb"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

type GRPCHandler struct {
	service Service
}

func NewGRPCHandler(service Service) *GRPCHandler {
	return &GRPCHandler{service: service}
}

func (handler *GRPCHandler) CreateEncounter(ctx context.Context, req *pb.CreateEncounterRequest) (*pb.CreateEncounterResponse, error) {
	input := CreateEncounterInput{
		PatientFHIRID:  req.PatientFhirId,
		PractitionerID: req.PractitionerId,
		ReasonCode:     req.ReasonCode,
		ReasonDisplay:  req.ReasonDisplay,
	}

	createdEncounter, err := handler.service.CreateEncounter(ctx, input)
	if err != nil {
		return nil, apperrors.ToGRPCStatus(err)
	}

	return &pb.CreateEncounterResponse{EncounterFhirId: createdEncounter.FHIRResourceID}, nil
}

func (handler *GRPCHandler) GetEncounter(ctx context.Context, req *pb.GetEncounterRequest) (*pb.GetEncounterResponse, error) {
	if req.EncounterFhirId == "" {
		return nil, apperrors.ErrBadRequest.WithFields(map[string]string{"encounter_fhir_id": "is required"})
	}

	encounter, err := handler.service.GetEncounter(ctx, req.EncounterFhirId)
	if err != nil {
		return nil, apperrors.ToGRPCStatus(err)
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
		return nil, apperrors.ToGRPCStatus(err)
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
