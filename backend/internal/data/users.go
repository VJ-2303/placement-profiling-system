package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID              int64     `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	Name            string    `json:"name"`
	OfficialEmail   string    `json:"official_email"`
	ProfileImageURL string    `json:"profile_image_url,omitempty"`
	Version         int       `json:"-"`
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(user *User) error {

	query := `INSERT INTO students (name,official_email,profile_image_url)
			  VALUES($1,$2,$3)
			  RETURNING id, created_at, version
	`

	args := []any{user.Name, user.OfficialEmail, user.ProfileImageURL}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
}

func (m UserModel) GetByEmail(email string) (*User, error) {

	query := `
		SELECT  id, created_at, name, official_email, profile_image_url, version
		FROM students
		WHERE official_email = $1
	`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.OfficialEmail,
		&user.ProfileImageURL,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}
