package models

import (
	"database/sql"
)

type Models struct {
	Students           StudentModel
	StudentDetails     StudentDetailsModel
	StudentParents     StudentParentsModel
	StudentAcademics   StudentAcademicsModel
	StudentAspirations StudentAspirationsModel
	Skills             SkillsModel
	StudentSkills      StudentSkillsModel
	Admins             AdminModel
	DB                 *sql.DB
}

func NewModels(db *sql.DB) Models {
	return Models{
		Students:           StudentModel{DB: db},
		StudentDetails:     StudentDetailsModel{DB: db},
		StudentParents:     StudentParentsModel{DB: db},
		StudentAcademics:   StudentAcademicsModel{DB: db},
		StudentAspirations: StudentAspirationsModel{DB: db},
		Skills:             SkillsModel{DB: db},
		StudentSkills:      StudentSkillsModel{DB: db},
		Admins:             AdminModel{DB: db},
		DB:                 db,
	}
}
