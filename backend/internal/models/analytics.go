package models

import (
	"context"
	"database/sql"
	"time"
)

// DashboardStats contains main dashboard statistics
type DashboardStats struct {
	TotalStudents        int     `json:"total_students"`
	ProfilesCompleted    int     `json:"profiles_completed"`
	ProfileCompletionPct float64 `json:"profile_completion_pct"`
	StudentsPlaced       int     `json:"students_placed"`
	StudentsNotPlaced    int     `json:"students_not_placed"`
	StudentsInProcess    int     `json:"students_in_process"`
	HigherStudies        int     `json:"higher_studies"`
	PlacementPct         float64 `json:"placement_pct"`
	AvgPackage           float64 `json:"avg_package"`
	MaxPackage           float64 `json:"max_package"`
	MinPackage           float64 `json:"min_package"`
	TotalCompanies       int     `json:"total_companies"`
}

// BatchStats contains statistics per batch
type BatchStats struct {
	BatchYear         int     `json:"batch_year"`
	TotalStudents     int     `json:"total_students"`
	ProfilesCompleted int     `json:"profiles_completed"`
	StudentsPlaced    int     `json:"students_placed"`
	PlacementPct      float64 `json:"placement_pct"`
	AvgPackage        float64 `json:"avg_package"`
	MaxPackage        float64 `json:"max_package"`
}

// SkillStats contains statistics about skills
type SkillStats struct {
	SkillName    string `json:"skill_name"`
	Category     string `json:"category"`
	Beginners    int    `json:"beginners"`
	Intermediate int    `json:"intermediate"`
	Advanced     int    `json:"advanced"`
	Experts      int    `json:"experts"`
	Total        int    `json:"total"`
}

// CGPADistribution contains CGPA range distribution
type CGPADistribution struct {
	Range     string `json:"range"`
	Count     int    `json:"count"`
	Placed    int    `json:"placed"`
	NotPlaced int    `json:"not_placed"`
}

// CompanyStats contains placement statistics by company
type CompanyStats struct {
	CompanyName string  `json:"company_name"`
	HiredCount  int     `json:"hired_count"`
	AvgPackage  float64 `json:"avg_package"`
	MaxPackage  float64 `json:"max_package"`
}

type AnalyticsModel struct {
	DB *sql.DB
}

// GetDashboardStats retrieves main dashboard statistics
func (m AnalyticsModel) GetDashboardStats(batchYear *int) (*DashboardStats, error) {
	var stats DashboardStats
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Main student stats query
	studentQuery := `
		SELECT 
			COUNT(*),
			COUNT(*) FILTER (WHERE is_profile_completed = true),
			COUNT(*) FILTER (WHERE placement_status = 'placed'),
			COUNT(*) FILTER (WHERE placement_status = 'not_placed'),
			COUNT(*) FILTER (WHERE placement_status = 'in_process'),
			COUNT(*) FILTER (WHERE placement_status = 'higher_studies')
		FROM students s
		LEFT JOIN batches b ON s.batch_id = b.id
		WHERE ($1::int IS NULL OR b.year = $1)`

	err := m.DB.QueryRowContext(ctx, studentQuery, batchYear).Scan(
		&stats.TotalStudents,
		&stats.ProfilesCompleted,
		&stats.StudentsPlaced,
		&stats.StudentsNotPlaced,
		&stats.StudentsInProcess,
		&stats.HigherStudies,
	)
	if err != nil {
		return nil, err
	}

	// Calculate percentages
	if stats.TotalStudents > 0 {
		stats.ProfileCompletionPct = float64(stats.ProfilesCompleted) / float64(stats.TotalStudents) * 100
		eligibleForPlacement := stats.TotalStudents - stats.HigherStudies
		if eligibleForPlacement > 0 {
			stats.PlacementPct = float64(stats.StudentsPlaced) / float64(eligibleForPlacement) * 100
		}
	}

	// Package stats
	packageQuery := `
		SELECT 
			COALESCE(AVG(p.package_lpa), 0),
			COALESCE(MAX(p.package_lpa), 0),
			COALESCE(MIN(p.package_lpa) FILTER (WHERE p.package_lpa > 0), 0)
		FROM placements p
		JOIN students s ON p.student_id = s.id
		LEFT JOIN batches b ON s.batch_id = b.id
		WHERE p.is_accepted = true AND ($1::int IS NULL OR b.year = $1)`

	err = m.DB.QueryRowContext(ctx, packageQuery, batchYear).Scan(
		&stats.AvgPackage,
		&stats.MaxPackage,
		&stats.MinPackage,
	)
	if err != nil {
		return nil, err
	}

	// Company count
	companyQuery := `SELECT COUNT(*) FROM companies WHERE is_active = true`
	err = m.DB.QueryRowContext(ctx, companyQuery).Scan(&stats.TotalCompanies)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetBatchWiseStats retrieves statistics per batch
func (m AnalyticsModel) GetBatchWiseStats() ([]BatchStats, error) {
	query := `
		SELECT 
			b.year,
			COUNT(DISTINCT s.id),
			COUNT(DISTINCT s.id) FILTER (WHERE s.is_profile_completed = true),
			COUNT(DISTINCT s.id) FILTER (WHERE s.placement_status = 'placed'),
			CASE 
				WHEN COUNT(DISTINCT s.id) > 0 
				THEN ROUND(COUNT(DISTINCT s.id) FILTER (WHERE s.placement_status = 'placed')::numeric / COUNT(DISTINCT s.id) * 100, 2)
				ELSE 0 
			END,
			COALESCE(AVG(p.package_lpa) FILTER (WHERE p.is_accepted = true), 0),
			COALESCE(MAX(p.package_lpa) FILTER (WHERE p.is_accepted = true), 0)
		FROM batches b
		LEFT JOIN students s ON s.batch_id = b.id
		LEFT JOIN placements p ON s.id = p.student_id
		WHERE b.is_active = true
		GROUP BY b.year
		ORDER BY b.year DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []BatchStats
	for rows.Next() {
		var s BatchStats
		if err := rows.Scan(
			&s.BatchYear, &s.TotalStudents, &s.ProfilesCompleted,
			&s.StudentsPlaced, &s.PlacementPct, &s.AvgPackage, &s.MaxPackage,
		); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}

	return stats, rows.Err()
}

// GetSkillStats retrieves skill-wise statistics
func (m AnalyticsModel) GetSkillStats() ([]SkillStats, error) {
	query := `
		SELECT 
			sk.name,
			sk.category,
			COUNT(*) FILTER (WHERE ss.proficiency = 'beginner'),
			COUNT(*) FILTER (WHERE ss.proficiency = 'intermediate'),
			COUNT(*) FILTER (WHERE ss.proficiency = 'advanced'),
			COUNT(*) FILTER (WHERE ss.proficiency = 'expert'),
			COUNT(*)
		FROM skills sk
		LEFT JOIN student_skills ss ON sk.id = ss.skill_id
		WHERE sk.is_active = true
		GROUP BY sk.id, sk.name, sk.category
		HAVING COUNT(*) > 0
		ORDER BY COUNT(*) DESC
		LIMIT 30`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []SkillStats
	for rows.Next() {
		var s SkillStats
		if err := rows.Scan(
			&s.SkillName, &s.Category, &s.Beginners, &s.Intermediate,
			&s.Advanced, &s.Experts, &s.Total,
		); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}

	return stats, rows.Err()
}

// GetCGPADistribution retrieves CGPA distribution
func (m AnalyticsModel) GetCGPADistribution(batchYear *int) ([]CGPADistribution, error) {
	query := `
		SELECT 
			CASE 
				WHEN sa.cgpa_overall >= 9 THEN '9.0 - 10.0'
				WHEN sa.cgpa_overall >= 8 THEN '8.0 - 8.99'
				WHEN sa.cgpa_overall >= 7 THEN '7.0 - 7.99'
				WHEN sa.cgpa_overall >= 6 THEN '6.0 - 6.99'
				ELSE 'Below 6.0'
			END as range,
			COUNT(*),
			COUNT(*) FILTER (WHERE s.placement_status = 'placed'),
			COUNT(*) FILTER (WHERE s.placement_status = 'not_placed')
		FROM students s
		JOIN student_academics sa ON s.id = sa.student_id
		LEFT JOIN batches b ON s.batch_id = b.id
		WHERE sa.cgpa_overall IS NOT NULL AND ($1::int IS NULL OR b.year = $1)
		GROUP BY 1
		ORDER BY 
			CASE 
				WHEN range = '9.0 - 10.0' THEN 1
				WHEN range = '8.0 - 8.99' THEN 2
				WHEN range = '7.0 - 7.99' THEN 3
				WHEN range = '6.0 - 6.99' THEN 4
				ELSE 5
			END`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, batchYear)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []CGPADistribution
	for rows.Next() {
		var s CGPADistribution
		if err := rows.Scan(&s.Range, &s.Count, &s.Placed, &s.NotPlaced); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}

	return stats, rows.Err()
}

// GetCompanyStats retrieves placement statistics by company
func (m AnalyticsModel) GetCompanyStats(batchYear *int) ([]CompanyStats, error) {
	query := `
		SELECT 
			COALESCE(c.name, p.company_name),
			COUNT(*),
			COALESCE(AVG(p.package_lpa), 0),
			COALESCE(MAX(p.package_lpa), 0)
		FROM placements p
		LEFT JOIN companies c ON p.company_id = c.id
		JOIN students s ON p.student_id = s.id
		LEFT JOIN batches b ON s.batch_id = b.id
		WHERE p.is_accepted = true AND ($1::int IS NULL OR b.year = $1)
		GROUP BY COALESCE(c.name, p.company_name)
		ORDER BY COUNT(*) DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, batchYear)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []CompanyStats
	for rows.Next() {
		var s CompanyStats
		if err := rows.Scan(&s.CompanyName, &s.HiredCount, &s.AvgPackage, &s.MaxPackage); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}

	return stats, rows.Err()
}

// GetRecentActivity retrieves recent activity for dashboard
type RecentActivity struct {
	Type        string    `json:"type"` // registration, profile_completed, placed
	StudentName string    `json:"student_name"`
	Details     string    `json:"details"`
	Timestamp   time.Time `json:"timestamp"`
}

func (m AnalyticsModel) GetRecentActivity(limit int) ([]RecentActivity, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	// This is a simplified version - in production you'd use activity_logs table
	query := `
		(SELECT 'registration' as type, name, 'New student registered' as details, created_at
		 FROM students ORDER BY created_at DESC LIMIT $1)
		UNION ALL
		(SELECT 'placed' as type, s.name, CONCAT('Placed at ', COALESCE(c.name, p.company_name)) as details, p.created_at
		 FROM placements p
		 JOIN students s ON p.student_id = s.id
		 LEFT JOIN companies c ON p.company_id = c.id
		 WHERE p.is_accepted = true
		 ORDER BY p.created_at DESC LIMIT $1)
		ORDER BY created_at DESC
		LIMIT $1`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []RecentActivity
	for rows.Next() {
		var a RecentActivity
		if err := rows.Scan(&a.Type, &a.StudentName, &a.Details, &a.Timestamp); err != nil {
			return nil, err
		}
		activities = append(activities, a)
	}

	return activities, rows.Err()
}

// GetBatches retrieves all batches
type Batch struct {
	ID       int  `json:"id"`
	Year     int  `json:"year"`
	IsActive bool `json:"is_active"`
}

func (m AnalyticsModel) GetBatches() ([]Batch, error) {
	query := `SELECT id, year, is_active FROM batches ORDER BY year DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var batches []Batch
	for rows.Next() {
		var b Batch
		if err := rows.Scan(&b.ID, &b.Year, &b.IsActive); err != nil {
			return nil, err
		}
		batches = append(batches, b)
	}

	return batches, rows.Err()
}
