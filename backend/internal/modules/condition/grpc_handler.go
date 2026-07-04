package condition

import (
	"context"
	"errors"

	pb "github.com/healthcare/backend/internal/modules/condition/pb"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

type GRPCHandler struct {
	service Service
}

func NewGRPCHandler(service Service) *GRPCHandler {
	return &GRPCHandler{service: service}
}

func mapConditionError(err error) error {
	if errors.Is(err, ErrConditionNotFound) {
		return apperrors.ErrConditionNotFound.ToGRPC()
	}
	return apperrors.ToGRPCStatus(err)
}

func (handler *GRPCHandler) CreateCondition(ctx context.Context, req *pb.CreateConditionRequest) (*pb.CreateConditionResponse, error) {
	violations := make(map[string]string)
	if req.PatientFhirId == "" {
		violations["patient_fhir_id"] = "is required"
	}
	if req.Icd10Code == "" {
		violations["icd10_code"] = "is required"
	}
	if len(violations) > 0 {
		return nil, apperrors.ErrBadRequest.WithFields(violations)
	}

	condition := &Condition{
		EncounterFHIRID: req.EncounterFhirId,
		PatientFHIRID:   req.PatientFhirId,
		ICD10Code:       req.Icd10Code,
		CodeDisplay:     req.CodeDisplay,
		ClinicalStatus:  req.ClinicalStatus,
	}

	createdCondition, err := handler.service.CreateCondition(ctx, condition)
	if err != nil {
		return nil, mapConditionError(err)
	}

	return &pb.CreateConditionResponse{ConditionFhirId: createdCondition.FHIRResourceID}, nil
}

func (handler *GRPCHandler) GetConditions(ctx context.Context, req *pb.GetConditionsRequest) (*pb.GetConditionsResponse, error) {
	if req.PatientFhirId == "" {
		return nil, apperrors.ErrBadRequest.WithFields(map[string]string{"patient_fhir_id": "is required"})
	}

	conditions, err := handler.service.GetConditionsByPatient(ctx, req.PatientFhirId)
	if err != nil {
		return nil, mapConditionError(err)
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
