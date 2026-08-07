package allergy

import (
	"context"

	pb "github.com/healthcare/backend/internal/modules/allergy/pb"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

type GRPCHandler struct {
	service Service
}

func NewGRPCHandler(service Service) *GRPCHandler {
	return &GRPCHandler{service: service}
}

func (handler *GRPCHandler) CreateAllergyIntolerance(ctx context.Context, req *pb.CreateAllergyIntoleranceRequest) (*pb.CreateAllergyIntoleranceResponse, error) {
	input := CreateAllergyInput{
		PatientFHIRID:   req.PatientFhirId,
		AllergenCode:    req.AllergenCode,
		AllergenDisplay: req.AllergenDisplay,
		ClinicalStatus:  req.ClinicalStatus,
		Reaction:        req.Reaction,
	}

	createdAllergy, err := handler.service.CreateAllergyIntolerance(ctx, input)
	if err != nil {
		return nil, apperrors.ToGRPCStatus(err)
	}

	return &pb.CreateAllergyIntoleranceResponse{AllergyFhirId: createdAllergy.FHIRResourceID}, nil
}

func (handler *GRPCHandler) GetAllergyIntolerances(ctx context.Context, req *pb.GetAllergyIntolerancesRequest) (*pb.GetAllergyIntolerancesResponse, error) {
	if req.PatientFhirId == "" {
		return nil, apperrors.ErrBadRequest.WithFields(map[string]string{"patient_fhir_id": "is required"})
	}

	allergies, err := handler.service.GetAllergyIntolerancesByPatient(ctx, req.PatientFhirId)
	if err != nil {
		return nil, apperrors.ToGRPCStatus(err)
	}

	pbAllergies := make([]*pb.AllergyIntolerance, 0, len(allergies))
	for _, allergy := range allergies {
		pbAllergies = append(pbAllergies, &pb.AllergyIntolerance{
			FhirId:          allergy.FHIRResourceID,
			AllergenDisplay: allergy.AllergenDisplay,
			ClinicalStatus:  allergy.ClinicalStatus,
			Reaction:        allergy.Reaction,
		})
	}

	return &pb.GetAllergyIntolerancesResponse{Allergies: pbAllergies}, nil
}
