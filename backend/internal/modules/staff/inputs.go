package staff

type CreateEmployeeInput struct {
	CreatedBy    string
	FullName     string
	Email        string
	Role         string
	CRMNumber    string
	DepartmentID string
}

type SetEmployeeActiveInput struct {
	EmployeeID string
	IsActive   bool
}
