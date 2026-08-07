package condition

import (
	"context"

	pb "github.com/healthcare/backend/internal/modules/condition/pb"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

type GRPCHandler struct {
	service Service
}

func NewGRPCHandler(service Service) *GRPCHandler {
	return &GRPCHandler{service: service}
}

func (handler *GRPCHandler) CreateCondition(ctx context.Context, req *pb.CreateConditionRequest) (*pb.CreateConditionResponse, error) {
	input := CreateConditionInput{
		PatientFHIRID:   req.PatientFhirId,
		EncounterFHIRID: req.EncounterFhirId,
		ICD10Code:       req.Icd10Code,
		CodeDisplay:     req.CodeDisplay,
		ClinicalStatus:  req.ClinicalStatus,
	}

	createdCondition, err := handler.service.CreateCondition(ctx, input)
	if err != nil {
		return nil, apperrors.ToGRPCStatus(err)
	}

	return &pb.CreateConditionResponse{ConditionFhirId: createdCondition.FHIRResourceID}, nil
}

func (handler *GRPCHandler) GetConditions(ctx context.Context, req *pb.GetConditionsRequest) (*pb.GetConditionsResponse, error) {
	if req.PatientFhirId == "" {
		return nil, apperrors.ErrBadRequest.WithFields(map[string]string{"patient_fhir_id": "is required"})
	}

	conditions, err := handler.service.GetConditionsByPatient(ctx, req.PatientFhirId)
	if err != nil {
		return nil, apperrors.ToGRPCStatus(err)
	}

	pbConditions := make([]*pb.Condition, 0, len(conditions))
	for _, condition := range conditions {
		pbConditions = append(pbConditions, &pb.Condition{
			FhirId:         condition.FHIRResourceID,
			Icd10Code:      condition.ICD10Code,
			CodeDisplay:    condition.CodeDisplay,
			ClinicalStatus: condition.ClinicalStatus,
		})
	}

	return &pb.GetConditionsResponse{Conditions: pbConditions}, nil
}
