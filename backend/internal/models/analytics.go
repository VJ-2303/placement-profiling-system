package models

import (
	"context"
	"database/sql"
	"time"
)

type Analytics struct {
	TotalStudents    int `json:"total_students"`
	ProfileCompleted int `json:"profile_completed"`
}

type AnalyticsModel struct {
	DB *sql.DB
}

func (m AnalyticsModel) GetDashboardAnalytics() (*Analytics, error) {
	query := `SELECT
				COUNT(*),
				COUNT(*) FILTER (WHERE is_profile_completed = TRUE)
			FROM
				students`

	var analytics Analytics

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query).Scan(&analytics.TotalStudents, &analytics.ProfileCompleted)

	if err != nil {
		return nil, err
	}
	return &analytics, nil
}
