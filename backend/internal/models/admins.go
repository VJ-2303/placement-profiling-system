package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Admin struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	OfficialEmail string    `json:"email"`
	CreatedAt     time.Time `json:"created_at"`
}

type AdminModel struct {
	DB *sql.DB
}

func (m AdminModel) Insert(admin *Admin) error {

	query := `
		INSERT INTO admins (name, email)
		VALUES ($1,$2)
		RETURNING id, created_at
	`
	args := []any{admin.Name, admin.OfficialEmail}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&admin.ID, admin.CreatedAt,
	)
}

func (m AdminModel) GetByEmail(email string) (*Admin, error) {

	if email == "" {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, name , email, created_at from admins where email = $1
	`

	var admin Admin

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&admin.ID,
		&admin.Name,
		&admin.OfficialEmail,
		&admin.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &admin, nil
}

func (m AdminModel) GetByID(id int64) (*Admin, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, name , email, created_at from admins where id = $1
	`

	var admin Admin

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&admin.ID,
		&admin.Name,
		&admin.OfficialEmail,
		&admin.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &admin, nil
}
