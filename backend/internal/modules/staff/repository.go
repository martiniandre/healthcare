package staff

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateEmployee(ctx context.Context, employee *Employee) error
	GetEmployeeByID(ctx context.Context, employeeID uuid.UUID) (*Employee, error)
	GetEmployeeByUserID(ctx context.Context, userID uuid.UUID) (*Employee, error)
	ListEmployees(ctx context.Context, search string, role string) ([]*Employee, error)
	DeactivateEmployee(ctx context.Context, employeeID uuid.UUID) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (staffRepository *repository) CreateEmployee(ctx context.Context, employee *Employee) error {
	query := `INSERT INTO employees (id, user_id, full_name, email, role, crm_number, is_active, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := staffRepository.db.Exec(ctx, query,
		employee.ID, employee.UserID, employee.FullName, employee.Email,
		employee.Role, employee.CRMNumber, employee.IsActive, employee.CreatedAt, employee.UpdatedAt,
	)
	return err
}

func (staffRepository *repository) GetEmployeeByID(ctx context.Context, employeeID uuid.UUID) (*Employee, error) {
	query := `SELECT id, user_id, full_name, email, role, crm_number, is_active, created_at, updated_at
			  FROM employees WHERE id = $1`

	employee := &Employee{}
	err := staffRepository.db.QueryRow(ctx, query, employeeID).Scan(
		&employee.ID, &employee.UserID, &employee.FullName, &employee.Email,
		&employee.Role, &employee.CRMNumber, &employee.IsActive, &employee.CreatedAt, &employee.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return employee, nil
}

func (staffRepository *repository) GetEmployeeByUserID(ctx context.Context, userID uuid.UUID) (*Employee, error) {
	query := `SELECT id, user_id, full_name, email, role, crm_number, is_active, created_at, updated_at
			  FROM employees WHERE user_id = $1`

	employee := &Employee{}
	err := staffRepository.db.QueryRow(ctx, query, userID).Scan(
		&employee.ID, &employee.UserID, &employee.FullName, &employee.Email,
		&employee.Role, &employee.CRMNumber, &employee.IsActive, &employee.CreatedAt, &employee.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return employee, nil
}

func (staffRepository *repository) ListEmployees(ctx context.Context, search string, role string) ([]*Employee, error) {
	query := `SELECT id, user_id, full_name, email, role, crm_number, is_active, created_at, updated_at
			  FROM employees WHERE is_active = true`
	
	args := []interface{}{}
	argId := 1

	if role != "" && role != "All" {
		query += fmt.Sprintf(" AND role = $%d", argId)
		args = append(args, role)
		argId++
	}

	if search != "" {
		query += fmt.Sprintf(" AND (full_name ILIKE $%d OR email ILIKE $%d)", argId, argId)
		args = append(args, "%"+search+"%")
		argId++
	}

	query += ` ORDER BY full_name ASC`

	rows, err := staffRepository.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employees := make([]*Employee, 0)
	for rows.Next() {
		employee := &Employee{}
		err := rows.Scan(
			&employee.ID, &employee.UserID, &employee.FullName, &employee.Email,
			&employee.Role, &employee.CRMNumber, &employee.IsActive, &employee.CreatedAt, &employee.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

func (staffRepository *repository) DeactivateEmployee(ctx context.Context, employeeID uuid.UUID) error {
	query := `UPDATE employees SET is_active = false, updated_at = NOW() WHERE id = $1`
	_, err := staffRepository.db.Exec(ctx, query, employeeID)
	return err
}
