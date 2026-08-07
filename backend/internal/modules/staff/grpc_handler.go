package staff

import (
	"context"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/staff/pb"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

type GRPCHandler struct {
	service Service
}

func NewGRPCHandler(service Service) *GRPCHandler {
	return &GRPCHandler{service: service}
}

func optionalStringValue(stringValuePointer *string) string {
	if stringValuePointer == nil {
		return ""
	}
	return *stringValuePointer
}

func mapEmployeeToResponse(employee *Employee) *pb.GetEmployeeResponse {
	return &pb.GetEmployeeResponse{
		EmployeeId:     employee.ID.String(),
		FullName:       employee.FullName,
		Email:          employee.Email,
		Role:           string(employee.Role),
		CrmNumber:      optionalStringValue(employee.CRMNumber),
		FhirResourceId: optionalStringValue(employee.FHIRResourceID),
		IsActive:       employee.IsActive,
	}
}

func (handler *GRPCHandler) CreateEmployee(ctx context.Context, req *pb.CreateEmployeeRequest) (*pb.CreateEmployeeResponse, error) {
	input := CreateEmployeeInput{
		CreatedBy: req.CreatedBy,
		FullName:  req.FullName,
		Email:     req.Email,
		Role:      req.Role,
		CRMNumber: req.CrmNumber,
	}

	employee, err := handler.service.CreateEmployee(ctx, input)
	if err != nil {
		return nil, apperrors.ToGRPCStatus(err)
	}

	return &pb.CreateEmployeeResponse{
		EmployeeId:     employee.ID.String(),
		FhirResourceId: optionalStringValue(employee.FHIRResourceID),
	}, nil
}

func (handler *GRPCHandler) GetEmployee(ctx context.Context, req *pb.GetEmployeeRequest) (*pb.GetEmployeeResponse, error) {
	employeeID, err := uuid.Parse(req.EmployeeId)
	if err != nil {
		return nil, apperrors.ErrBadRequest.ToGRPC()
	}

	employee, err := handler.service.GetEmployee(ctx, employeeID)
	if err != nil {
		return nil, apperrors.ToGRPCStatus(err)
	}

	return mapEmployeeToResponse(employee), nil
}

func (handler *GRPCHandler) ListEmployees(ctx context.Context, req *pb.ListEmployeesRequest) (*pb.ListEmployeesResponse, error) {
	employees, listError := handler.service.ListEmployees(ctx, "", "")
	if listError != nil {
		return nil, apperrors.ToGRPCStatus(listError)
	}

	employeeResponses := make([]*pb.GetEmployeeResponse, 0, len(employees))
	for _, employee := range employees {
		employeeResponses = append(employeeResponses, mapEmployeeToResponse(employee))
	}

	return &pb.ListEmployeesResponse{Employees: employeeResponses}, nil
}

func (handler *GRPCHandler) DeactivateEmployee(ctx context.Context, req *pb.DeactivateEmployeeRequest) (*pb.DeactivateEmployeeResponse, error) {
	employeeID, err := uuid.Parse(req.EmployeeId)
	if err != nil {
		return nil, apperrors.ErrBadRequest.ToGRPC()
	}

	err = handler.service.DeactivateEmployee(ctx, employeeID)
	if err != nil {
		return nil, apperrors.ToGRPCStatus(err)
	}

	return &pb.DeactivateEmployeeResponse{}, nil
}
