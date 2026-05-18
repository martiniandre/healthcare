package staff

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/staff/pb"
	"github.com/healthcare/backend/internal/shared/apperrors"
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

func (handler *GRPCHandler) CreateEmployee(ctx context.Context, req *pb.CreateEmployeeRequest) (*pb.CreateEmployeeResponse, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, apperrors.ErrBadRequest.ToGRPC()
	}

	employee, err := handler.service.CreateEmployee(ctx, userID, req.FullName, req.Email, req.Role, req.CrmNumber)
	if err != nil {
		return nil, mapStaffError(err)
	}

	return &pb.CreateEmployeeResponse{
		EmployeeId: employee.ID.String(),
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
		EmployeeId: employee.ID.String(),
		FullName:   employee.FullName,
		Email:      employee.Email,
		Role:       string(employee.Role),
		CrmNumber:  employee.CRMNumber,
		IsActive:   employee.IsActive,
	}, nil
}

func (handler *GRPCHandler) ListEmployees(ctx context.Context, req *pb.ListEmployeesRequest) (*pb.ListEmployeesResponse, error) {
	employees, err := handler.service.ListEmployees(ctx)
	if err != nil {
		return nil, mapStaffError(err)
	}

	employeeResponses := make([]*pb.GetEmployeeResponse, 0, len(employees))
	for _, employee := range employees {
		employeeResponses = append(employeeResponses, &pb.GetEmployeeResponse{
			EmployeeId: employee.ID.String(),
			FullName:   employee.FullName,
			Email:      employee.Email,
			Role:       string(employee.Role),
			CrmNumber:  employee.CRMNumber,
			IsActive:   employee.IsActive,
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
