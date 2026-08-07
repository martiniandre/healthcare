package patients

import (
	"context"

	"github.com/healthcare/backend/internal/modules/patients/pb"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

type GRPCHandler struct {
	service Service
}

func NewGRPCHandler(service Service) *GRPCHandler {
	return &GRPCHandler{service: service}
}

func mapPatientToResponse(patient *Patient) *pb.GetPatientResponse {
	return &pb.GetPatientResponse{
		PatientId:      patient.ID.String(),
		FhirResourceId: patient.FHIRResourceID,
		FullName:       patient.FullName,
		BirthDate:      patient.BirthDate.Format("2006-01-02"),
		DocumentId:     patient.DocumentID,
		PhoneNumber:    patient.PhoneNumber,
	}
}

func (handler *GRPCHandler) CreatePatient(ctx context.Context, req *pb.CreatePatientRequest) (*pb.CreatePatientResponse, error) {
	input := CreatePatientInput{
		FullName:    req.FullName,
		BirthDate:   req.BirthDate,
		DocumentID:  req.DocumentID,
		PhoneNumber: req.PhoneNumber,
	}

	patient, err := handler.service.CreatePatient(ctx, input)
	if err != nil {
		return nil, apperrors.ToGRPCStatus(err)
	}

	return &pb.CreatePatientResponse{
		PatientId:      patient.ID.String(),
		FhirResourceId: patient.FHIRResourceID,
	}, nil
}

func (handler *GRPCHandler) GetPatient(ctx context.Context, req *pb.GetPatientRequest) (*pb.GetPatientResponse, error) {
	patient, err := handler.service.GetPatient(ctx, req.FhirResourceId)
	if err != nil {
		return nil, apperrors.ToGRPCStatus(err)
	}

	return mapPatientToResponse(patient), nil
}

func (handler *GRPCHandler) GetPatientByDocument(ctx context.Context, req *pb.GetPatientByDocumentRequest) (*pb.GetPatientResponse, error) {
	patient, err := handler.service.GetPatientByDocument(ctx, req.DocumentId)
	if err != nil {
		return nil, apperrors.ToGRPCStatus(err)
	}

	return mapPatientToResponse(patient), nil
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
		return nil, apperrors.ToGRPCStatus(listError)
	}

	patientResponses := make([]*pb.GetPatientResponse, 0, len(patientsList))
	for _, patient := range patientsList {
		patientResponses = append(patientResponses, mapPatientToResponse(patient))
	}

	return &pb.ListPatientsResponse{Patients: patientResponses}, nil
}
