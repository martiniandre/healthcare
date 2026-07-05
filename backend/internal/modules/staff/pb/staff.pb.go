package pb

import "context"

type StaffServiceServer interface {
	CreateEmployee(ctx context.Context, req *CreateEmployeeRequest) (*CreateEmployeeResponse, error)
	GetEmployee(ctx context.Context, req *GetEmployeeRequest) (*GetEmployeeResponse, error)
	ListEmployees(ctx context.Context, req *ListEmployeesRequest) (*ListEmployeesResponse, error)
	DeactivateEmployee(ctx context.Context, req *DeactivateEmployeeRequest) (*DeactivateEmployeeResponse, error)
}

type CreateEmployeeRequest struct {
	CreatedBy  string
	FullName   string
	Email      string
	Role       string
	CrmNumber  string
}

type CreateEmployeeResponse struct {
	EmployeeId      string
	FhirResourceId  string
}

type GetEmployeeRequest struct {
	EmployeeId string
}

type GetEmployeeResponse struct {
	EmployeeId      string
	FullName        string
	Email           string
	Role            string
	CrmNumber       string
	FhirResourceId  string
	IsActive        bool
}

type ListEmployeesRequest struct{}

type ListEmployeesResponse struct {
	Employees []*GetEmployeeResponse
}

type DeactivateEmployeeRequest struct {
	EmployeeId string
}

type DeactivateEmployeeResponse struct{}

func RegisterStaffServiceServer(_ interface{}, server StaffServiceServer) {}
