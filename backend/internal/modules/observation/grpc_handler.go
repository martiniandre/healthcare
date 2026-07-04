package observation

import (
	"context"
	"errors"

	pb "github.com/healthcare/backend/internal/modules/observation/pb"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

type GRPCHandler struct {
	service Service
}

func NewGRPCHandler(service Service) *GRPCHandler {
	return &GRPCHandler{service: service}
}

func mapObservationError(err error) error {
	if errors.Is(err, ErrObservationNotFound) {
		return apperrors.ErrObservationNotFound.ToGRPC()
	}
	return apperrors.ToGRPCStatus(err)
}

func (handler *GRPCHandler) CreateObservation(ctx context.Context, req *pb.CreateObservationRequest) (*pb.CreateObservationResponse, error) {
	violations := make(map[string]string)
	if req.EncounterFhirId == "" {
		violations["encounter_fhir_id"] = "is required"
	}
	if req.PatientFhirId == "" {
		violations["patient_fhir_id"] = "is required"
	}
	if req.LoincCode == "" {
		violations["loinc_code"] = "is required"
	}
	if len(violations) > 0 {
		return nil, apperrors.ErrBadRequest.WithFields(violations)
	}

	observation := &Observation{
		EncounterFHIRID: req.EncounterFhirId,
		PatientFHIRID:   req.PatientFhirId,
		LoincCode:       req.LoincCode,
		CodeDisplay:     req.CodeDisplay,
		ValueQuantity:   req.ValueQuantity,
		ValueUnit:       req.ValueUnit,
	}

	createdObservation, err := handler.service.CreateObservation(ctx, observation)
	if err != nil {
		return nil, mapObservationError(err)
	}

	return &pb.CreateObservationResponse{ObservationFhirId: createdObservation.FHIRResourceID}, nil
}

func (handler *GRPCHandler) GetObservations(ctx context.Context, req *pb.GetObservationsRequest) (*pb.GetObservationsResponse, error) {
	if req.EncounterFhirId == "" {
		return nil, apperrors.ErrBadRequest.WithFields(map[string]string{"encounter_fhir_id": "is required"})
	}

	observations, err := handler.service.GetObservationsByEncounter(ctx, req.EncounterFhirId)
	if err != nil {
		return nil, mapObservationError(err)
	}

	pbObservations := make([]*pb.Observation, 0, len(observations))
	for _, observation := range observations {
		pbObservations = append(pbObservations, &pb.Observation{
			FhirId:        observation.FHIRResourceID,
			LoincCode:     observation.LoincCode,
			CodeDisplay:   observation.CodeDisplay,
			ValueQuantity: observation.ValueQuantity,
			ValueUnit:     observation.ValueUnit,
		})
	}

	return &pb.GetObservationsResponse{Observations: pbObservations}, nil
}
