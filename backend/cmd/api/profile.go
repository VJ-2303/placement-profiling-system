package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/VJ-2303/placement-profiling-system/internal/models"
)

func (app *application) createStudentProfileHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := app.authenticateStudent(r)
	app.logger.Print(claims)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}
	studentID := int64(claims.UserID)

	var input FlatProfileRequest
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	tx, err := app.models.DB.Begin()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE students SET roll_no = $1, name = $2 , photo = $3 WHERE id = $4",
		input.RollNo, input.Name, input.Photo, studentID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	dob, err := time.Parse("2006-01-02", input.DateOfBirth)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("invalid date format for date_of_birth: %v", err))
		return
	}

	pincode := ""
	if input.Pincode != nil {
		pincode = *input.Pincode
	}

	details := models.StudentDetails{
		StudentID:             studentID,
		DateOfBirth:           dob,
		MobileNumber:          input.MobileNumber,
		AlternateMobileNumber: input.AltMobileNumber,
		PersonalEmail:         input.PersonalEmail,
		LinkedinProfile:       input.LinkedInUrl,
		Address:               input.Address,
		City:                  input.City,
		Pincode:               pincode,
		AdhaarNo:              input.AdhaarNo,
		ResidenceType:         input.ResidenceType,
		Strength:              input.Strength,
		Weakness:              input.Weakness,
		Remarks:               input.Remarks,
	}
	err = app.models.StudentDetails.Insert(tx, &details)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	parents := models.StudentParents{
		StudentID:            studentID,
		FatherName:           input.FatherName,
		FatherMobile:         input.FatherMobile,
		FatherOccupation:     input.FatherOccupation,
		FatherCompanyDetails: input.FatherCompanyDetails,
		FatherEmail:          input.FatherEmail,
		MotherName:           input.MotherName,
		MotherMobile:         input.MotherMobile,
		MotherOccupation:     input.MotherOccupation,
		MotherEmail:          input.MotherEmail,
	}
	err = app.models.StudentParents.Insert(tx, &parents)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	academics := models.StudentAcademics{
		StudentID:         studentID,
		TenthPercentage:   input.TenthPercentage,
		TwelthPercentage:  input.TwelthPercentage,
		CgpaSem1:          input.CgpaSem1,
		CgpaSem2:          input.CgpaSem2,
		CgpaSem3:          input.CgpaSem3,
		CgpaSem4:          input.CgpaSem4,
		CgpaOverall:       input.CgpaOverall,
		CurrentBacklogs:   input.CurrentBacklogs,
		HasBacklogHistory: input.HasBacklogHistory,
	}
	err = app.models.StudentAcademics.Insert(tx, &academics)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	aspirations := models.StudentAspirations{
		StudentID:           studentID,
		CompanyAim:          input.CompanyAim,
		TargetPackage:       input.TargetPackage,
		Certifications:      input.Certifications,
		Awards:              input.Awards,
		Workshops:           input.Workshops,
		Internships:         input.Internships,
		HackathonsAttended:  input.HackathonsAttended,
		Extracurriculars:    input.Extracurriculars,
		ClubParticipation:   input.ClubParticipation,
		FuturePath:          input.FuturePath,
		CommunicationSkills: input.CommunicationSkills,
	}
	err = app.models.StudentAspirations.Insert(tx, &aspirations)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	skillMap, err := app.models.Skills.GetAllAsMap(tx)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	skillsToInsert := map[string]string{
		"C":                                input.SkillC,
		"C++":                              input.SkillCpp,
		"JAVA":                             input.SkillJava,
		"PYTHON":                           input.SkillPython,
		"Node.js":                          input.SkillNodeJs,
		"SQL Database":                     input.SkillSql,
		"NoSQL Database":                   input.SkillNoSql,
		"Web Developement":                 input.SkillWebDev,
		"PHP":                              input.SkillPhp,
		"Mobile App development-flutter":   input.SkillFlutter,
		"Aptitude level":                   input.SkillAptitude,
		"logical and verbal Reasoning":     input.SkillReasoning,
		"DataStructure":                    input.ConceptDataStructures,
		"DBMS":                             input.ConceptDbms,
		"OOPS":                             input.ConceptOops,
		"Problem Solving/Coding Tests":     input.ConceptProblemSolving,
		"Computer Networks":                input.ConceptNetworks,
		"Operating System":                 input.ConceptOs,
		"Design and Analysis of Algorithm": input.ConceptAlgos,
		"Git/Github":                       input.ToolGit,
		"Linux/Unix":                       input.ToolLinux,
		"Cloud Basics (AWS/Azure/GCP)":     input.ToolCloud,
		"Hacker Rank":                      input.ToolHackerRank,
		"Hacker Earth":                     input.ToolHackerEarth,
	}
	for skillName, proficiency := range skillsToInsert {
		if skillID, ok := skillMap[skillName]; ok && proficiency != "" {
			err = app.models.StudentSkills.Insert(tx, studentID, skillID, proficiency)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
		}
	}

	err = app.models.Students.SetStudentProfileCompleted(tx, studentID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := tx.Commit(); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"message": "Profile created successfully"}, nil)
}

func (app *application) getStudentProfileHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := app.authenticateStudent(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}
	studentID := int64(claims.UserID)

	profile, err := app.models.Students.GetFullProfile(studentID)
	if err != nil {
		switch err {
		case models.ErrRecordNotFound:
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"profile": profile}, nil)
}
