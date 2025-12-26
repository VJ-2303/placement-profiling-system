package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// Skill represents a master skill entry
type Skill struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	Category     SkillCategory `json:"category"`
	Description  *string       `json:"description"`
	IsActive     bool          `json:"is_active"`
	DisplayOrder int           `json:"display_order"`
}

type SkillModel struct {
	DB *sql.DB
}

// GetAll retrieves all active skills
func (m SkillModel) GetAll() ([]Skill, error) {
	query := `
		SELECT id, name, category, description, is_active, display_order
		FROM skills
		WHERE is_active = true
		ORDER BY category, display_order`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []Skill
	for rows.Next() {
		var s Skill
		if err := rows.Scan(&s.ID, &s.Name, &s.Category, &s.Description, &s.IsActive, &s.DisplayOrder); err != nil {
			return nil, err
		}
		skills = append(skills, s)
	}

	return skills, rows.Err()
}

// GetByCategory retrieves skills by category
func (m SkillModel) GetByCategory(category SkillCategory) ([]Skill, error) {
	query := `
		SELECT id, name, category, description, is_active, display_order
		FROM skills
		WHERE category = $1 AND is_active = true
		ORDER BY display_order`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []Skill
	for rows.Next() {
		var s Skill
		if err := rows.Scan(&s.ID, &s.Name, &s.Category, &s.Description, &s.IsActive, &s.DisplayOrder); err != nil {
			return nil, err
		}
		skills = append(skills, s)
	}

	return skills, rows.Err()
}

// GetGroupedByCategory retrieves skills grouped by category
func (m SkillModel) GetGroupedByCategory() (map[SkillCategory][]Skill, error) {
	skills, err := m.GetAll()
	if err != nil {
		return nil, err
	}

	grouped := make(map[SkillCategory][]Skill)
	for _, skill := range skills {
		grouped[skill.Category] = append(grouped[skill.Category], skill)
	}

	return grouped, nil
}

// GetByID retrieves a skill by ID
func (m SkillModel) GetByID(id int) (*Skill, error) {
	query := `
		SELECT id, name, category, description, is_active, display_order
		FROM skills
		WHERE id = $1`

	var s Skill
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.Name, &s.Category, &s.Description, &s.IsActive, &s.DisplayOrder,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &s, nil
}

// Insert adds a new skill
func (m SkillModel) Insert(skill *Skill) error {
	query := `
		INSERT INTO skills (name, category, description, display_order)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query,
		skill.Name, skill.Category, skill.Description, skill.DisplayOrder,
	).Scan(&skill.ID)
}

// GetAsMap returns a map of skill name to ID
func (m SkillModel) GetAsMap() (map[string]int, error) {
	skills, err := m.GetAll()
	if err != nil {
		return nil, err
	}

	skillMap := make(map[string]int)
	for _, s := range skills {
		skillMap[s.Name] = s.ID
	}

	return skillMap, nil
}
