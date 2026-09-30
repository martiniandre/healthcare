package staff

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateEmployee(ctx context.Context, employee *Employee) error
	GetEmployeeByID(ctx context.Context, employeeID uuid.UUID) (*Employee, error)
	ListEmployees(ctx context.Context, search string, role string) ([]*Employee, error)
	SetEmployeeActive(ctx context.Context, employeeID uuid.UUID, isActive bool) error
	GetDepartmentByID(ctx context.Context, departmentID uuid.UUID) (*Department, error)
	GetDefaultDepartment(ctx context.Context) (*Department, error)
	ListDepartments(ctx context.Context) ([]*Department, error)
	GetRoleIDByCode(ctx context.Context, roleCode string) (uuid.UUID, error)
	UpdateEmployeeFHIRResourceID(ctx context.Context, employeeID uuid.UUID, fhirResourceID string) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (staffRepository *repository) CreateEmployee(ctx context.Context, employee *Employee) error {
	query := `INSERT INTO employees (id, full_name, email, role_id, crm_number, fhir_resource_id, department_id, created_by, is_active, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := staffRepository.db.Exec(ctx, query,
		employee.ID, employee.FullName, employee.Email,
		employee.RoleID, employee.CRMNumber, employee.FHIRResourceID, employee.DepartmentID, employee.CreatedBy, employee.IsActive, employee.CreatedAt, employee.UpdatedAt,
	)
	if err != nil {
		var postgresError *pgconn.PgError
		if errors.As(err, &postgresError) && postgresError.Code == "23505" {
			return apperrors.ErrEmployeeAlreadyExists
		}
		if errors.As(err, &postgresError) && postgresError.Code == "23503" {
			return apperrors.InvalidArgument("invalid employee input", map[string]string{"department_id": "unknown department"})
		}
	}
	return err
}

func (staffRepository *repository) UpdateEmployeeFHIRResourceID(ctx context.Context, employeeID uuid.UUID, fhirResourceID string) error {
	query := `UPDATE employees SET fhir_resource_id = $1, updated_at = NOW() WHERE id = $2`
	_, err := staffRepository.db.Exec(ctx, query, fhirResourceID, employeeID)
	return err
}

func (staffRepository *repository) GetEmployeeByID(ctx context.Context, employeeID uuid.UUID) (*Employee, error) {
	query := `SELECT e.id, e.full_name, e.email, r.code, e.crm_number, e.fhir_resource_id, e.department_id, d.name, e.created_by, e.is_active, e.created_at, e.updated_at
			  FROM employees e
			  JOIN departments d ON d.id = e.department_id
			  JOIN roles r ON r.id = e.role_id
			  WHERE e.id = $1`

	employee := &Employee{}
	err := staffRepository.db.QueryRow(ctx, query, employeeID).Scan(
		&employee.ID, &employee.FullName, &employee.Email,
		&employee.Role, &employee.CRMNumber, &employee.FHIRResourceID, &employee.DepartmentID, &employee.DepartmentName,
		&employee.CreatedBy, &employee.IsActive, &employee.CreatedAt, &employee.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("failed to get employee: %w", apperrors.ErrEmployeeNotFound)
		}
		return nil, err
	}
	return employee, nil
}

func (staffRepository *repository) ListEmployees(ctx context.Context, search string, role string) ([]*Employee, error) {
	query := `SELECT e.id, e.full_name, e.email, r.code, e.crm_number, e.fhir_resource_id, e.department_id, d.name, e.created_by, e.is_active, e.created_at, e.updated_at
			  FROM employees e
			  JOIN departments d ON d.id = e.department_id
			  JOIN roles r ON r.id = e.role_id`

	args := []interface{}{}
	argId := 1

	if role != "" && role != "All" {
		query += fmt.Sprintf(" WHERE r.code = $%d", argId)
		args = append(args, role)
		argId++
	}

	if search != "" {
		searchClause := fmt.Sprintf("(e.full_name ILIKE $%d OR e.email ILIKE $%d)", argId, argId)
		if role != "" && role != "All" {
			query += fmt.Sprintf(" AND %s", searchClause)
		} else {
			query += fmt.Sprintf(" WHERE %s", searchClause)
		}
		args = append(args, "%"+search+"%")
		argId++
	}

	query += ` ORDER BY e.is_active DESC, e.full_name ASC`

	rows, err := staffRepository.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employees := make([]*Employee, 0)
	for rows.Next() {
		employee := &Employee{}
		err := rows.Scan(
			&employee.ID, &employee.FullName, &employee.Email,
			&employee.Role, &employee.CRMNumber, &employee.FHIRResourceID, &employee.DepartmentID, &employee.DepartmentName,
			&employee.CreatedBy, &employee.IsActive, &employee.CreatedAt, &employee.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, rows.Err()
}

func (staffRepository *repository) SetEmployeeActive(ctx context.Context, employeeID uuid.UUID, isActive bool) error {
	query := `UPDATE employees SET is_active = $1, updated_at = NOW() WHERE id = $2`
	commandTag, err := staffRepository.db.Exec(ctx, query, isActive, employeeID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("failed to update employee status: %w", apperrors.ErrEmployeeNotFound)
	}
	return nil
}

func (staffRepository *repository) GetDepartmentByID(ctx context.Context, departmentID uuid.UUID) (*Department, error) {
	query := `SELECT id, name, sort_order, is_active, created_at FROM departments WHERE id = $1`

	department := &Department{}
	err := staffRepository.db.QueryRow(ctx, query, departmentID).Scan(
		&department.ID, &department.Name, &department.SortOrder, &department.IsActive, &department.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("failed to get department: %w", apperrors.ErrDepartmentNotFound)
		}
		return nil, err
	}
	return department, nil
}

func (staffRepository *repository) GetDefaultDepartment(ctx context.Context) (*Department, error) {
	query := `SELECT id, name, sort_order, is_active, created_at FROM departments
			  WHERE is_active = true ORDER BY sort_order ASC, name ASC LIMIT 1`

	department := &Department{}
	err := staffRepository.db.QueryRow(ctx, query).Scan(
		&department.ID, &department.Name, &department.SortOrder, &department.IsActive, &department.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("failed to get default department: %w", apperrors.ErrDepartmentNotFound)
		}
		return nil, err
	}
	return department, nil
}

func (staffRepository *repository) ListDepartments(ctx context.Context) ([]*Department, error) {
	query := `SELECT id, name, sort_order, is_active, created_at FROM departments WHERE is_active = true ORDER BY sort_order ASC, name ASC`

	rows, err := staffRepository.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	departments := make([]*Department, 0)
	for rows.Next() {
		department := &Department{}
		err := rows.Scan(&department.ID, &department.Name, &department.SortOrder, &department.IsActive, &department.CreatedAt)
		if err != nil {
			return nil, err
		}
		departments = append(departments, department)
	}
	return departments, rows.Err()
}

func (staffRepository *repository) GetRoleIDByCode(ctx context.Context, roleCode string) (uuid.UUID, error) {
	var roleID uuid.UUID
	err := staffRepository.db.QueryRow(ctx, `SELECT id FROM roles WHERE code = $1`, roleCode).Scan(&roleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("failed to get role: %w", apperrors.ErrRoleNotFound)
		}
		return uuid.Nil, err
	}
	return roleID, nil
}
