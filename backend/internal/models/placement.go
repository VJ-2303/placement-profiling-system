package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// PlacementRecord represents a student's placement
type PlacementRecord struct {
	ID             int64      `json:"id"`
	StudentID      int64      `json:"student_id"`
	CompanyID      *int64     `json:"company_id"`
	CompanyName    string     `json:"company_name"`
	JobRole        *string    `json:"job_role"`
	PackageLPA     *float64   `json:"package_lpa"`
	PackageCTC     *string    `json:"package_ctc"`
	JoiningDate    *string    `json:"joining_date"`
	OfferDate      *string    `json:"offer_date"`
	OfferType      *string    `json:"offer_type"` // full_time, internship, ppo
	JobLocation    *string    `json:"job_location"`
	IsAccepted     bool       `json:"is_accepted"`
	VerifiedBy     *int64     `json:"verified_by"`
	VerifiedAt     *time.Time `json:"verified_at"`
	OfferLetterURL *string    `json:"offer_letter_url"`
	Remarks        *string    `json:"remarks"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type PlacementModel struct {
	DB *sql.DB
}

// GetByStudentID retrieves placement record for a student
func (m PlacementModel) GetByStudentID(studentID int64) (*PlacementRecord, error) {
	query := `
		SELECT id, student_id, company_id, 
		       COALESCE(company_name, (SELECT name FROM companies WHERE id = company_id)),
		       job_role, package_lpa, package_ctc, joining_date::text, offer_date::text,
		       offer_type, job_location, is_accepted, verified_by, verified_at,
		       offer_letter_url, remarks, created_at, updated_at
		FROM placements
		WHERE student_id = $1 AND is_accepted = true
		ORDER BY created_at DESC
		LIMIT 1`

	var p PlacementRecord
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, studentID).Scan(
		&p.ID, &p.StudentID, &p.CompanyID, &p.CompanyName, &p.JobRole, &p.PackageLPA,
		&p.PackageCTC, &p.JoiningDate, &p.OfferDate, &p.OfferType, &p.JobLocation,
		&p.IsAccepted, &p.VerifiedBy, &p.VerifiedAt, &p.OfferLetterURL, &p.Remarks,
		&p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

// Insert creates a new placement record
func (m PlacementModel) Insert(p *PlacementRecord) error {
	query := `
		INSERT INTO placements (
			student_id, company_id, company_name, job_role, package_lpa, package_ctc,
			joining_date, offer_date, offer_type, job_location, is_accepted,
			verified_by, offer_letter_url, remarks
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var joiningDate, offerDate interface{}
	if p.JoiningDate != nil && *p.JoiningDate != "" {
		joiningDate = *p.JoiningDate
	}
	if p.OfferDate != nil && *p.OfferDate != "" {
		offerDate = *p.OfferDate
	}

	return m.DB.QueryRowContext(ctx, query,
		p.StudentID, p.CompanyID, p.CompanyName, p.JobRole, p.PackageLPA, p.PackageCTC,
		joiningDate, offerDate, p.OfferType, p.JobLocation, p.IsAccepted,
		p.VerifiedBy, p.OfferLetterURL, p.Remarks,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

// Update updates a placement record
func (m PlacementModel) Update(p *PlacementRecord) error {
	query := `
		UPDATE placements
		SET company_id = $1, company_name = $2, job_role = $3, package_lpa = $4,
		    package_ctc = $5, joining_date = $6, offer_date = $7, offer_type = $8,
		    job_location = $9, is_accepted = $10, offer_letter_url = $11, remarks = $12
		WHERE id = $13
		RETURNING updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var joiningDate, offerDate interface{}
	if p.JoiningDate != nil && *p.JoiningDate != "" {
		joiningDate = *p.JoiningDate
	}
	if p.OfferDate != nil && *p.OfferDate != "" {
		offerDate = *p.OfferDate
	}

	return m.DB.QueryRowContext(ctx, query,
		p.CompanyID, p.CompanyName, p.JobRole, p.PackageLPA, p.PackageCTC,
		joiningDate, offerDate, p.OfferType, p.JobLocation, p.IsAccepted,
		p.OfferLetterURL, p.Remarks, p.ID,
	).Scan(&p.UpdatedAt)
}

// Delete deletes a placement record
func (m PlacementModel) Delete(id int64) error {
	query := `DELETE FROM placements WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// Verify marks a placement as verified by admin
func (m PlacementModel) Verify(id int64, adminID int64) error {
	query := `
		UPDATE placements
		SET verified_by = $1, verified_at = NOW()
		WHERE id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, adminID, id)
	return err
}

// GetAllPlacements retrieves all placement records with student info
type PlacementWithStudent struct {
	PlacementRecord
	StudentName     string  `json:"student_name"`
	StudentRoll     *string `json:"student_roll"`
	StudentEmail    string  `json:"student_email"`
	StudentPhotoURL *string `json:"student_photo_url"`
	BatchYear       *int    `json:"batch_year"`
}

func (m PlacementModel) GetAll() ([]PlacementWithStudent, error) {
	query := `
		SELECT p.id, p.student_id, p.company_id,
		       COALESCE(p.company_name, c.name), p.job_role, p.package_lpa,
		       p.package_ctc, p.joining_date::text, p.offer_date::text, p.offer_type,
		       p.job_location, p.is_accepted, p.verified_by, p.verified_at,
		       p.offer_letter_url, p.remarks, p.created_at, p.updated_at,
		       s.name, s.roll_no, s.official_email, s.photo_url, b.year
		FROM placements p
		JOIN students s ON p.student_id = s.id
		LEFT JOIN companies c ON p.company_id = c.id
		LEFT JOIN batches b ON s.batch_id = b.id
		WHERE p.is_accepted = true
		ORDER BY p.created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var placements []PlacementWithStudent
	for rows.Next() {
		var p PlacementWithStudent
		if err := rows.Scan(
			&p.ID, &p.StudentID, &p.CompanyID, &p.CompanyName, &p.JobRole, &p.PackageLPA,
			&p.PackageCTC, &p.JoiningDate, &p.OfferDate, &p.OfferType, &p.JobLocation,
			&p.IsAccepted, &p.VerifiedBy, &p.VerifiedAt, &p.OfferLetterURL, &p.Remarks,
			&p.CreatedAt, &p.UpdatedAt, &p.StudentName, &p.StudentRoll, &p.StudentEmail, &p.StudentPhotoURL, &p.BatchYear,
		); err != nil {
			return nil, err
		}
		placements = append(placements, p)
	}

	return placements, rows.Err()
}
