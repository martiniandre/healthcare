package observation

import (
	"context"

	pb "github.com/healthcare/backend/internal/modules/observation/pb"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

type GRPCHandler struct {
	service Service
}

func NewGRPCHandler(service Service) *GRPCHandler {
	return &GRPCHandler{service: service}
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

	input := CreateObservationInput{
		EncounterFHIRID: req.EncounterFhirId,
		PatientFHIRID:   req.PatientFhirId,
		LoincCode:       req.LoincCode,
		CodeDisplay:     req.CodeDisplay,
		ValueQuantity:   req.ValueQuantity,
		ValueUnit:       req.ValueUnit,
	}

	createdObservation, err := handler.service.CreateObservation(ctx, input)
	if err != nil {
		return nil, apperrors.ToGRPCStatus(err)
	}

	return &pb.CreateObservationResponse{ObservationFhirId: createdObservation.FHIRResourceID}, nil
}

func (handler *GRPCHandler) CreateObservationBatch(ctx context.Context, req *pb.CreateObservationBatchRequest) (*pb.CreateObservationBatchResponse, error) {
	violations := make(map[string]string)
	if req.EncounterFhirId == "" {
		violations["encounter_fhir_id"] = "is required"
	}
	if req.PatientFhirId == "" {
		violations["patient_fhir_id"] = "is required"
	}
	if len(violations) > 0 {
		return nil, apperrors.ErrBadRequest.WithFields(violations)
	}

	input := CreateObservationBatchInput{
		EncounterFHIRID: req.EncounterFhirId,
		PatientFHIRID:   req.PatientFhirId,
	}
	if req.Panel != nil {
		input.HeartRate = req.Panel.HeartRate
		input.BodyTemperature = req.Panel.BodyTemperature
		input.SystolicBloodPressure = req.Panel.SystolicBloodPressure
		input.DiastolicBloodPressure = req.Panel.DiastolicBloodPressure
		input.OxygenSaturation = req.Panel.OxygenSaturation
		input.RespiratoryRate = req.Panel.RespiratoryRate
		input.WeightKilograms = req.Panel.WeightKilograms
		input.HeightCentimeters = req.Panel.HeightCentimeters
	}

	createdObservations, batchErr := handler.service.CreateObservationBatch(ctx, input)
	if batchErr != nil {
		return nil, apperrors.ToGRPCStatus(batchErr)
	}

	pbObservations := make([]*pb.Observation, 0, len(createdObservations))
	for _, observation := range createdObservations {
		pbObservations = append(pbObservations, &pb.Observation{
			FhirId:        observation.FHIRResourceID,
			LoincCode:     observation.LoincCode,
			CodeDisplay:   observation.CodeDisplay,
			ValueQuantity: observation.ValueQuantity,
			ValueUnit:     observation.ValueUnit,
			NotPerformed:  observation.NotPerformed,
		})
	}

	return &pb.CreateObservationBatchResponse{Observations: pbObservations}, nil
}

func (handler *GRPCHandler) GetObservations(ctx context.Context, req *pb.GetObservationsRequest) (*pb.GetObservationsResponse, error) {
	if req.EncounterFhirId == "" {
		return nil, apperrors.ErrBadRequest.WithFields(map[string]string{"encounter_fhir_id": "is required"})
	}

	observations, err := handler.service.GetObservationsByEncounter(ctx, req.EncounterFhirId)
	if err != nil {
		return nil, apperrors.ToGRPCStatus(err)
	}

	pbObservations := make([]*pb.Observation, 0, len(observations))
	for _, observation := range observations {
		pbObservations = append(pbObservations, &pb.Observation{
			FhirId:        observation.FHIRResourceID,
			LoincCode:     observation.LoincCode,
			CodeDisplay:   observation.CodeDisplay,
			ValueQuantity: observation.ValueQuantity,
			ValueUnit:     observation.ValueUnit,
			NotPerformed:  observation.NotPerformed,
		})
	}

	return &pb.GetObservationsResponse{Observations: pbObservations}, nil
}
