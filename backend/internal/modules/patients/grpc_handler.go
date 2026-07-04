package patients

import (
	"context"
	"errors"
	"strings"

	"github.com/healthcare/backend/internal/modules/patients/pb"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/validator"
)

func mapPatientError(err error) error {
	switch {
	case errors.Is(err, ErrPatientNotFound):
		return apperrors.ErrPatientNotFound.ToGRPC()
	case errors.Is(err, ErrPatientAlreadyExists):
		return apperrors.ErrPatientAlreadyExists.ToGRPC()
	case err.Error() == "invalid birth date format, expected YYYY-MM-DD" || err.Error() == "birth date must be in the past":
		return apperrors.ErrBadRequest.ToGRPC()
	default:
		return apperrors.ToGRPCStatus(err)
	}
}

type GRPCHandler struct {
	service Service
}

func NewGRPCHandler(service Service) *GRPCHandler {
	return &GRPCHandler{service: service}
}

func (handler *GRPCHandler) CreatePatient(ctx context.Context, req *pb.CreatePatientRequest) (*pb.CreatePatientResponse, error) {
	violations := make(map[string]string)
	if strings.TrimSpace(req.FullName) == "" {
		violations["full_name"] = "full name is required"
	}
	if strings.TrimSpace(req.BirthDate) == "" {
		violations["birth_date"] = "birth date is required"
	}
	if strings.TrimSpace(req.DocumentID) == "" || !validator.IsValidCPF(req.DocumentID) {
		violations["document_id"] = "invalid CPF format"
	}
	if strings.TrimSpace(req.PhoneNumber) == "" || !validator.IsValidPhone(req.PhoneNumber) {
		violations["phone_number"] = "invalid phone format"
	}
	if len(violations) > 0 {
		return nil, apperrors.ErrBadRequest.WithFields(violations)
	}

	patient, err := handler.service.CreatePatient(ctx, req.FullName, req.BirthDate, req.DocumentID, req.PhoneNumber)
	if err != nil {
		return nil, mapPatientError(err)
	}

	return &pb.CreatePatientResponse{
		PatientId:      patient.ID.String(),
		FhirResourceId: patient.FHIRResourceID,
	}, nil
}

func (handler *GRPCHandler) GetPatient(ctx context.Context, req *pb.GetPatientRequest) (*pb.GetPatientResponse, error) {
	patient, err := handler.service.GetPatient(ctx, req.FhirResourceId)
	if err != nil {
		return nil, mapPatientError(err)
	}

	return &pb.GetPatientResponse{
		PatientId:      patient.ID.String(),
		FhirResourceId: patient.FHIRResourceID,
		FullName:       patient.FullName,
		BirthDate:      patient.BirthDate.Format("2006-01-02"),
		DocumentId:     patient.DocumentID,
		PhoneNumber:    patient.PhoneNumber,
	}, nil
}

func (handler *GRPCHandler) GetPatientByDocument(ctx context.Context, req *pb.GetPatientByDocumentRequest) (*pb.GetPatientResponse, error) {
	patient, err := handler.service.GetPatientByDocument(ctx, req.DocumentId)
	if err != nil {
		return nil, mapPatientError(err)
	}

	return &pb.GetPatientResponse{
		PatientId:      patient.ID.String(),
		FhirResourceId: patient.FHIRResourceID,
		FullName:       patient.FullName,
		BirthDate:      patient.BirthDate.Format("2006-01-02"),
		DocumentId:     patient.DocumentID,
		PhoneNumber:    patient.PhoneNumber,
	}, nil
}

func (handler *GRPCHandler) ListPatients(ctx context.Context, req *pb.ListPatientsRequest) (*pb.ListPatientsResponse, error) {
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 50
	}

	patientsList, listError := handler.service.ListPatients(ctx, req.Search, req.SortField, req.SortDirection, page, limit)
	if listError != nil {
		return nil, mapPatientError(listError)
	}

	patientResponses := make([]*pb.GetPatientResponse, 0, len(patientsList))
	for _, patient := range patientsList {
		patientResponses = append(patientResponses, &pb.GetPatientResponse{
			PatientId:      patient.ID.String(),
			FhirResourceId: patient.FHIRResourceID,
			FullName:       patient.FullName,
			BirthDate:      patient.BirthDate.Format("2006-01-02"),
			DocumentId:     patient.DocumentID,
			PhoneNumber:    patient.PhoneNumber,
		})
	}

	return &pb.ListPatientsResponse{Patients: patientResponses}, nil
}
