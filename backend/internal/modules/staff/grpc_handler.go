package staff

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/staff/pb"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/validator"
)

func mapStaffError(err error) error {
	if errors.Is(err, ErrEmployeeNotFound) {
		return apperrors.ErrEmployeeNotFound.ToGRPC()
	}
	return apperrors.ToGRPCStatus(err)
}

type GRPCHandler struct {
	service Service
}

func NewGRPCHandler(service Service) *GRPCHandler {
	return &GRPCHandler{service: service}
}

func fhirResourceIDToString(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func (handler *GRPCHandler) CreateEmployee(ctx context.Context, req *pb.CreateEmployeeRequest) (*pb.CreateEmployeeResponse, error) {
	violations := make(map[string]string)
	createdBy, err := uuid.Parse(req.CreatedBy)
	if err != nil {
		violations["created_by"] = "invalid UUID format"
	}
	if strings.TrimSpace(req.FullName) == "" {
		violations["full_name"] = "full name is required"
	}
	if strings.TrimSpace(req.Email) == "" || !validator.IsValidEmail(req.Email) {
		violations["email"] = "invalid email format"
	}
	if strings.TrimSpace(req.Role) == "" {
		violations["role"] = "role is required"
	}
	if len(violations) > 0 {
		return nil, apperrors.ErrBadRequest.WithFields(violations)
	}

	employee, err := handler.service.CreateEmployee(ctx, createdBy, req.FullName, req.Email, req.Role, req.CrmNumber)
	if err != nil {
		return nil, mapStaffError(err)
	}

	return &pb.CreateEmployeeResponse{
		EmployeeId:     employee.ID.String(),
		FhirResourceId: fhirResourceIDToString(employee.FHIRResourceID),
	}, nil
}

func (handler *GRPCHandler) GetEmployee(ctx context.Context, req *pb.GetEmployeeRequest) (*pb.GetEmployeeResponse, error) {
	employeeID, err := uuid.Parse(req.EmployeeId)
	if err != nil {
		return nil, apperrors.ErrBadRequest.ToGRPC()
	}

	employee, err := handler.service.GetEmployee(ctx, employeeID)
	if err != nil {
		return nil, mapStaffError(err)
	}

	return &pb.GetEmployeeResponse{
		EmployeeId:     employee.ID.String(),
		FullName:       employee.FullName,
		Email:          employee.Email,
		Role:           string(employee.Role),
		CrmNumber:      getStringValue(employee.CRMNumber),
		FhirResourceId: fhirResourceIDToString(employee.FHIRResourceID),
		IsActive:       employee.IsActive,
	}, nil
}

func (handler *GRPCHandler) ListEmployees(ctx context.Context, req *pb.ListEmployeesRequest) (*pb.ListEmployeesResponse, error) {
	employees, listError := handler.service.ListEmployees(ctx, "", "")
	if listError != nil {
		return nil, mapStaffError(listError)
	}

	employeeResponses := make([]*pb.GetEmployeeResponse, 0, len(employees))
	for _, employee := range employees {
		employeeResponses = append(employeeResponses, &pb.GetEmployeeResponse{
			EmployeeId:     employee.ID.String(),
			FullName:       employee.FullName,
			Email:          employee.Email,
			Role:           string(employee.Role),
			CrmNumber:      getStringValue(employee.CRMNumber),
			FhirResourceId: fhirResourceIDToString(employee.FHIRResourceID),
			IsActive:       employee.IsActive,
		})
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
		return nil, mapStaffError(err)
	}

	return &pb.DeactivateEmployeeResponse{}, nil
}

func getStringValue(stringValuePointer *string) string {
	if stringValuePointer == nil {
		return ""
	}
	return *stringValuePointer
}
