package models

import (
	"database/sql"
)

type Models struct {
	Students StudentModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Students: StudentModel{DB: db},
	}
}
