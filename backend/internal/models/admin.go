package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// Admin represents a pre-registered admin user
type Admin struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phone       *string   `json:"phone"`
	Designation string    `json:"designation"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AdminModel struct {
	DB *sql.DB
}

// GetByEmail retrieves an admin by email (must be pre-registered)
func (m AdminModel) GetByEmail(email string) (*Admin, error) {
	query := `
		SELECT id, name, email, phone, designation, is_active, created_at, updated_at
		FROM admins
		WHERE email = $1 AND is_active = true`

	var admin Admin
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&admin.ID, &admin.Name, &admin.Email, &admin.Phone,
		&admin.Designation, &admin.IsActive, &admin.CreatedAt, &admin.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &admin, nil
}

// GetByID retrieves an admin by ID
func (m AdminModel) GetByID(id int64) (*Admin, error) {
	query := `
		SELECT id, name, email, phone, designation, is_active, created_at, updated_at
		FROM admins
		WHERE id = $1`

	var admin Admin
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&admin.ID, &admin.Name, &admin.Email, &admin.Phone,
		&admin.Designation, &admin.IsActive, &admin.CreatedAt, &admin.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &admin, nil
}

// GetAll retrieves all admins
func (m AdminModel) GetAll() ([]Admin, error) {
	query := `
		SELECT id, name, email, phone, designation, is_active, created_at, updated_at
		FROM admins
		ORDER BY name ASC`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var admins []Admin
	for rows.Next() {
		var admin Admin
		if err := rows.Scan(
			&admin.ID, &admin.Name, &admin.Email, &admin.Phone,
			&admin.Designation, &admin.IsActive, &admin.CreatedAt, &admin.UpdatedAt,
		); err != nil {
			return nil, err
		}
		admins = append(admins, admin)
	}

	return admins, rows.Err()
}

// Insert creates a new admin (for initial setup)
func (m AdminModel) Insert(admin *Admin) error {
	query := `
		INSERT INTO admins (name, email, phone, designation)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query,
		admin.Name, admin.Email, admin.Phone, admin.Designation,
	).Scan(&admin.ID, &admin.CreatedAt, &admin.UpdatedAt)
}

// Update updates an admin's info
func (m AdminModel) Update(admin *Admin) error {
	query := `
		UPDATE admins
		SET name = $1, phone = $2, designation = $3, is_active = $4
		WHERE id = $5
		RETURNING updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query,
		admin.Name, admin.Phone, admin.Designation, admin.IsActive, admin.ID,
	).Scan(&admin.UpdatedAt)
}
