package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// Company represents a recruiting company
type Company struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Website      *string   `json:"website"`
	Industry     *string   `json:"industry"`
	CompanyType  *string   `json:"company_type"` // product, service, startup, mnc
	Description  *string   `json:"description"`
	LogoURL      *string   `json:"logo_url"`
	HRName       *string   `json:"hr_name"`
	HREmail      *string   `json:"hr_email"`
	HRPhone      *string   `json:"hr_phone"`
	Headquarters *string   `json:"headquarters"`
	Locations    *string   `json:"locations"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CompanyModel struct {
	DB *sql.DB
}

// GetAll retrieves all companies
func (m CompanyModel) GetAll() ([]Company, error) {
	query := `
		SELECT id, name, website, industry, company_type, description, logo_url,
		       hr_name, hr_email, hr_phone, headquarters, locations, is_active,
		       created_at, updated_at
		FROM companies
		WHERE is_active = true
		ORDER BY name ASC`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []Company
	for rows.Next() {
		var c Company
		if err := rows.Scan(
			&c.ID, &c.Name, &c.Website, &c.Industry, &c.CompanyType, &c.Description,
			&c.LogoURL, &c.HRName, &c.HREmail, &c.HRPhone, &c.Headquarters, &c.Locations,
			&c.IsActive, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		companies = append(companies, c)
	}

	return companies, rows.Err()
}

// GetByID retrieves a company by ID
func (m CompanyModel) GetByID(id int64) (*Company, error) {
	query := `
		SELECT id, name, website, industry, company_type, description, logo_url,
		       hr_name, hr_email, hr_phone, headquarters, locations, is_active,
		       created_at, updated_at
		FROM companies
		WHERE id = $1`

	var c Company
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Website, &c.Industry, &c.CompanyType, &c.Description,
		&c.LogoURL, &c.HRName, &c.HREmail, &c.HRPhone, &c.Headquarters, &c.Locations,
		&c.IsActive, &c.CreatedAt, &c.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &c, nil
}

// Insert creates a new company
func (m CompanyModel) Insert(c *Company) error {
	query := `
		INSERT INTO companies (name, website, industry, company_type, description, logo_url,
		                       hr_name, hr_email, hr_phone, headquarters, locations)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query,
		c.Name, c.Website, c.Industry, c.CompanyType, c.Description, c.LogoURL,
		c.HRName, c.HREmail, c.HRPhone, c.Headquarters, c.Locations,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

// Update updates a company
func (m CompanyModel) Update(c *Company) error {
	query := `
		UPDATE companies
		SET name = $1, website = $2, industry = $3, company_type = $4, description = $5,
		    logo_url = $6, hr_name = $7, hr_email = $8, hr_phone = $9, headquarters = $10,
		    locations = $11, is_active = $12
		WHERE id = $13
		RETURNING updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query,
		c.Name, c.Website, c.Industry, c.CompanyType, c.Description, c.LogoURL,
		c.HRName, c.HREmail, c.HRPhone, c.Headquarters, c.Locations, c.IsActive, c.ID,
	).Scan(&c.UpdatedAt)
}

// Search searches companies by name
func (m CompanyModel) Search(search string) ([]Company, error) {
	query := `
		SELECT id, name, website, industry, company_type, description, logo_url,
		       hr_name, hr_email, hr_phone, headquarters, locations, is_active,
		       created_at, updated_at
		FROM companies
		WHERE name ILIKE $1 AND is_active = true
		ORDER BY name ASC
		LIMIT 20`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, "%"+search+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []Company
	for rows.Next() {
		var c Company
		if err := rows.Scan(
			&c.ID, &c.Name, &c.Website, &c.Industry, &c.CompanyType, &c.Description,
			&c.LogoURL, &c.HRName, &c.HREmail, &c.HRPhone, &c.Headquarters, &c.Locations,
			&c.IsActive, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		companies = append(companies, c)
	}

	return companies, rows.Err()
}

// Delete deletes a company by ID
func (m CompanyModel) Delete(id int64) error {
	query := `DELETE FROM companies WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
