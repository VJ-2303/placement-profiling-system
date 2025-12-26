package models

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound  = errors.New("record not found")
	ErrEditConflict    = errors.New("edit conflict")
	ErrDuplicateEmail  = errors.New("duplicate email")
	ErrDuplicateRollNo = errors.New("duplicate roll number")
)

type Models struct {
	Students   StudentModel
	Admins     AdminModel
	Skills     SkillModel
	Companies  CompanyModel
	Placements PlacementModel
	Analytics  AnalyticsModel
	DB         *sql.DB
}

func NewModels(db *sql.DB) Models {
	return Models{
		Students:   StudentModel{DB: db},
		Admins:     AdminModel{DB: db},
		Skills:     SkillModel{DB: db},
		Companies:  CompanyModel{DB: db},
		Placements: PlacementModel{DB: db},
		Analytics:  AnalyticsModel{DB: db},
		DB:         db,
	}
}
