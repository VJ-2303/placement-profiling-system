package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Student struct {
	ID              int64     `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	Name            string    `json:"name"`
	OfficialEmail   string    `json:"official_email"`
	ProfileImageURL string    `json:"profile_image_url"`
	Version         int       `json:"version"`
}

type StudentModel struct {
	DB *sql.DB
}

func (m StudentModel) Insert(student *Student) error {
	query := `
		INSERT INTO students (name, official_email, profile_image_url)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, version`

	args := []any{student.Name, student.OfficialEmail, student.ProfileImageURL}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&student.ID,
		&student.CreatedAt,
		&student.Version,
	)
}

func (m StudentModel) GetByEmail(email string) (*Student, error) {
	if email == "" {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, name, official_email, profile_image_url, version
		FROM students
		WHERE official_email = $1`

	var student Student

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&student.ID,
		&student.CreatedAt,
		&student.Name,
		&student.OfficialEmail,
		&student.ProfileImageURL,
		&student.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &student, nil
}

func (m StudentModel) GetByID(id int64) (*Student, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, name, official_email, profile_image_url, version
		FROM students
		WHERE id = $1`

	var student Student

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&student.ID,
		&student.CreatedAt,
		&student.Name,
		&student.OfficialEmail,
		&student.ProfileImageURL,
		&student.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &student, nil
}
