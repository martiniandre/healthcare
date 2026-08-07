package medication

import (
	"context"

	pb "github.com/healthcare/backend/internal/modules/medication/pb"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

type GRPCHandler struct {
	service Service
}

func NewGRPCHandler(service Service) *GRPCHandler {
	return &GRPCHandler{service: service}
}

func (handler *GRPCHandler) CreateMedicationRequest(ctx context.Context, req *pb.CreateMedicationRequestRequest) (*pb.CreateMedicationRequestResponse, error) {
	input := CreateMedicationInput{
		EncounterFHIRID:    req.EncounterFhirId,
		PatientFHIRID:      req.PatientFhirId,
		PractitionerFHIRID: req.PractitionerFhirId,
		MedicationCode:     req.MedicationCode,
		MedicationName:     req.MedicationName,
		DosageInstructions: req.DosageInstructions,
	}

	createdMedication, err := handler.service.CreateMedicationRequest(ctx, input)
	if err != nil {
		return nil, apperrors.ToGRPCStatus(err)
	}

	return &pb.CreateMedicationRequestResponse{MedicationRequestFhirId: createdMedication.FHIRResourceID}, nil
}

func (handler *GRPCHandler) GetMedicationRequests(ctx context.Context, req *pb.GetMedicationRequestsRequest) (*pb.GetMedicationRequestsResponse, error) {
	if req.EncounterFhirId == "" {
		return nil, apperrors.ErrBadRequest.WithFields(map[string]string{"encounter_fhir_id": "is required"})
	}

	medications, err := handler.service.GetMedicationRequestsByEncounter(ctx, req.EncounterFhirId)
	if err != nil {
		return nil, apperrors.ToGRPCStatus(err)
	}

	pbMedications := make([]*pb.MedicationRequest, 0, len(medications))
	for _, medication := range medications {
		pbMedications = append(pbMedications, &pb.MedicationRequest{
			FhirId:             medication.FHIRResourceID,
			MedicationName:     medication.MedicationName,
			DosageInstructions: medication.DosageInstructions,
			Status:             medication.Status,
		})
	}

	return &pb.GetMedicationRequestsResponse{MedicationRequests: pbMedications}, nil
}
