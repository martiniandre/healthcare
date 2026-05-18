package patients

import (
	"context"
	"errors"

	"github.com/healthcare/backend/internal/modules/patients/pb"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

func mapPatientError(err error) error {
	switch {
	case errors.Is(err, ErrPatientNotFound):
		return apperrors.ErrPatientNotFound.ToGRPC()
	case errors.Is(err, ErrPatientAlreadyExists):
		return apperrors.ErrPatientAlreadyExists.ToGRPC()
	case err.Error() == "invalid birth date format, expected YYYY-MM-DD":
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
	patientList, err := handler.service.ListPatients(ctx)
	if err != nil {
		return nil, mapPatientError(err)
	}

	patientResponses := make([]*pb.GetPatientResponse, 0, len(patientList))
	for _, patient := range patientList {
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
