package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/VJ-2303/placement-profiling-system/internal/models"
)

// ============================================
// STUDENT PROFILE HANDLERS
// ============================================

// getStudentProfile returns the student's full profile
func (app *application) getStudentProfile(w http.ResponseWriter, r *http.Request) {
	claims, err := app.authenticateStudent(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	profile, err := app.models.Students.GetFullProfile(claims.UserID)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	// Get placement info if placed
	if profile.Student.PlacementStatus == models.PlacementStatusPlaced {
		placement, _ := app.models.Placements.GetByStudentID(claims.UserID)
		profile.Placement = placement
	}

	app.writeJSON(w, http.StatusOK, envelope{"profile": profile}, nil)
}

// updateStudentProfile updates basic student info
func (app *application) updateStudentProfile(w http.ResponseWriter, r *http.Request) {
	claims, err := app.authenticateStudent(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	var input struct {
		Name       string  `json:"name"`
		RollNo     *string `json:"roll_no"`
		RegisterNo *string `json:"register_no"`
		BatchID    *int    `json:"batch_id"`
		PhotoURL   *string `json:"photo_url"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	student, err := app.models.Students.GetByID(claims.UserID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if input.Name != "" {
		student.Name = input.Name
	}
	if input.RollNo != nil {
		student.RollNo = input.RollNo
	}
	if input.RegisterNo != nil {
		student.RegisterNo = input.RegisterNo
	}
	if input.BatchID != nil {
		student.BatchID = input.BatchID
	}
	if input.PhotoURL != nil {
		student.PhotoURL = input.PhotoURL
	}

	if err := app.models.Students.UpdateBasicInfo(student); err != nil {
		if errors.Is(err, models.ErrEditConflict) {
			app.errorResponse(w, r, http.StatusConflict, "record was modified by another request")
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"student": student}, nil)
}

// updatePersonalDetails updates student personal details
func (app *application) updatePersonalDetails(w http.ResponseWriter, r *http.Request) {
	claims, err := app.authenticateStudent(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	// Custom input struct that accepts both Student and PersonalDetails fields
	var input struct {
		// Student fields (will update Student table)
		Name      string  `json:"name"`
		RollNo    *string `json:"roll_no"`
		BatchYear *int    `json:"batch_year"`

		// Personal details fields (with frontend aliases)
		DateOfBirth     *string `json:"date_of_birth"`
		Gender          *string `json:"gender"`
		BloodGroup      *string `json:"blood_group"`
		MobileNumber    *string `json:"mobile_number"`
		AlternateMobile *string `json:"alt_mobile_number"` // Frontend uses alt_mobile_number
		PersonalEmail   *string `json:"personal_email"`
		LinkedinURL     *string `json:"linkedin_url"`
		GithubURL       *string `json:"github_url"`
		PortfolioURL    *string `json:"portfolio_url"`
		AadhaarNumber   *string `json:"aadhaar_no"` // Frontend uses aadhaar_no
		Address         *string `json:"address"`
		City            *string `json:"city"`
		State           *string `json:"state"`
		Pincode         *string `json:"pincode"`
		ResidenceType   *string `json:"residence_type"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	tx, err := app.models.DB.Begin()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer tx.Rollback()

	// Update Student basic info if provided
	student, err := app.models.Students.GetByID(claims.UserID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if input.Name != "" {
		student.Name = input.Name
	}
	if input.RollNo != nil {
		student.RollNo = input.RollNo
	}
	if input.BatchYear != nil {
		// Need to find batch ID from year
		batchID, err := app.models.Students.GetBatchIDByYear(*input.BatchYear)
		if err == nil && batchID > 0 {
			student.BatchID = &batchID
		}
	}

	if err := app.models.Students.UpdateBasicInfo(student); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Update personal details
	personalDetails := &models.StudentPersonalDetails{
		StudentID:       claims.UserID,
		DateOfBirth:     input.DateOfBirth,
		Gender:          input.Gender,
		BloodGroup:      input.BloodGroup,
		MobileNumber:    input.MobileNumber,
		AlternateMobile: input.AlternateMobile,
		PersonalEmail:   input.PersonalEmail,
		LinkedinURL:     input.LinkedinURL,
		GithubURL:       input.GithubURL,
		PortfolioURL:    input.PortfolioURL,
		AadhaarNumber:   input.AadhaarNumber,
		Address:         input.Address,
		City:            input.City,
		State:           input.State,
		Pincode:         input.Pincode,
		ResidenceType:   input.ResidenceType,
	}

	if err := app.models.Students.UpsertPersonalDetails(tx, personalDetails); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := tx.Commit(); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "Personal details saved"}, nil)
}

// updateFamilyDetails updates student family details
func (app *application) updateFamilyDetails(w http.ResponseWriter, r *http.Request) {
	claims, err := app.authenticateStudent(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	// Custom input struct with frontend field aliases
	var input struct {
		FatherName         *string `json:"father_name"`
		FatherMobile       *string `json:"father_mobile"`
		FatherEmail        *string `json:"father_email"`
		FatherOccupation   *string `json:"father_occupation"`
		FatherCompany      *string `json:"father_company_details"` // Frontend uses father_company_details
		FatherAnnualIncome *string `json:"annual_income"`          // Frontend uses annual_income
		MotherName         *string `json:"mother_name"`
		MotherMobile       *string `json:"mother_mobile"`
		MotherEmail        *string `json:"mother_email"`
		MotherOccupation   *string `json:"mother_occupation"`
		MotherCompany      *string `json:"mother_company"`
		GuardianName       *string `json:"guardian_name"`
		GuardianMobile     *string `json:"guardian_mobile"`
		GuardianRelation   *string `json:"guardian_relation"`
		ResidenceType      *string `json:"residence_type"` // This goes to personal details
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	tx, err := app.models.DB.Begin()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer tx.Rollback()

	familyDetails := &models.StudentFamilyDetails{
		StudentID:          claims.UserID,
		FatherName:         input.FatherName,
		FatherMobile:       input.FatherMobile,
		FatherEmail:        input.FatherEmail,
		FatherOccupation:   input.FatherOccupation,
		FatherCompany:      input.FatherCompany,
		FatherAnnualIncome: input.FatherAnnualIncome,
		MotherName:         input.MotherName,
		MotherMobile:       input.MotherMobile,
		MotherEmail:        input.MotherEmail,
		MotherOccupation:   input.MotherOccupation,
		MotherCompany:      input.MotherCompany,
		GuardianName:       input.GuardianName,
		GuardianMobile:     input.GuardianMobile,
		GuardianRelation:   input.GuardianRelation,
	}

	if err := app.models.Students.UpsertFamilyDetails(tx, familyDetails); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := tx.Commit(); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "Family details saved"}, nil)
}

// updateAcademics updates student academic details
func (app *application) updateAcademics(w http.ResponseWriter, r *http.Request) {
	claims, err := app.authenticateStudent(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	// Custom input struct with frontend field aliases
	var input struct {
		TenthPercentage   *float64 `json:"tenth_percentage"`
		TenthBoard        *string  `json:"tenth_board"`
		TenthYear         *int     `json:"tenth_year"`
		TenthSchool       *string  `json:"tenth_school"`
		TwelfthPercentage *float64 `json:"twelfth_percentage"`
		TwelfthBoard      *string  `json:"twelfth_board"`
		TwelfthYear       *int     `json:"twelfth_year"`
		TwelfthSchool     *string  `json:"twelfth_school"`
		HasDiploma        bool     `json:"has_diploma"`
		DiplomaPercentage *float64 `json:"diploma_percentage"`
		DiplomaBranch     *string  `json:"diploma_branch"`
		DiplomaCollege    *string  `json:"diploma_college"`
		CGPASem1          *float64 `json:"cgpa_sem1"`
		CGPASem2          *float64 `json:"cgpa_sem2"`
		CGPASem3          *float64 `json:"cgpa_sem3"`
		CGPASem4          *float64 `json:"cgpa_sem4"`
		CGPASem5          *float64 `json:"cgpa_sem5"`
		CGPASem6          *float64 `json:"cgpa_sem6"`
		CGPASem7          *float64 `json:"cgpa_sem7"`
		CGPASem8          *float64 `json:"cgpa_sem8"`
		CGPAOverall       *float64 `json:"cgpa_overall"`
		CurrentBacklogs   int      `json:"current_backlogs"`
		HistoryOfBacklogs bool     `json:"has_backlog_history"` // Frontend uses has_backlog_history
		BacklogDetails    *string  `json:"backlog_details"`
		HasGapYear        bool     `json:"has_gap_year"`
		GapYearReason     *string  `json:"gap_year_reason"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	tx, err := app.models.DB.Begin()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer tx.Rollback()

	academics := &models.StudentAcademics{
		StudentID:         claims.UserID,
		TenthPercentage:   input.TenthPercentage,
		TenthBoard:        input.TenthBoard,
		TenthYear:         input.TenthYear,
		TenthSchool:       input.TenthSchool,
		TwelfthPercentage: input.TwelfthPercentage,
		TwelfthBoard:      input.TwelfthBoard,
		TwelfthYear:       input.TwelfthYear,
		TwelfthSchool:     input.TwelfthSchool,
		HasDiploma:        input.HasDiploma,
		DiplomaPercentage: input.DiplomaPercentage,
		DiplomaBranch:     input.DiplomaBranch,
		DiplomaCollege:    input.DiplomaCollege,
		CGPASem1:          input.CGPASem1,
		CGPASem2:          input.CGPASem2,
		CGPASem3:          input.CGPASem3,
		CGPASem4:          input.CGPASem4,
		CGPASem5:          input.CGPASem5,
		CGPASem6:          input.CGPASem6,
		CGPASem7:          input.CGPASem7,
		CGPASem8:          input.CGPASem8,
		CGPAOverall:       input.CGPAOverall,
		CurrentBacklogs:   input.CurrentBacklogs,
		HistoryOfBacklogs: input.HistoryOfBacklogs,
		BacklogDetails:    input.BacklogDetails,
		HasGapYear:        input.HasGapYear,
		GapYearReason:     input.GapYearReason,
	}

	if err := app.models.Students.UpsertAcademics(tx, academics); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := tx.Commit(); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "Academic details saved"}, nil)
}

// updateAchievements updates student achievements
func (app *application) updateAchievements(w http.ResponseWriter, r *http.Request) {
	claims, err := app.authenticateStudent(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	// Custom input struct with frontend field aliases
	var input struct {
		Certifications         *string `json:"certifications"`
		Awards                 *string `json:"awards"`
		Workshops              *string `json:"workshops"`
		Internships            *string `json:"internships"`
		Projects               *string `json:"projects"`
		LeetcodeProfile        *string `json:"leetcode_profile"`
		HackerrankProfile      *string `json:"hackerrank_profile"`
		CodeforcesProfile      *string `json:"codeforces_profile"`
		CodechefProfile        *string `json:"codechef_profile"`
		LeetcodeRating         *int    `json:"leetcode_rating"`
		ProblemsSolved         *int    `json:"problems_solved"`
		HackathonsParticipated int     `json:"hackathons_participated"`
		HackathonsWon          int     `json:"hackathons_won"`
		HackathonDetails       *string `json:"hackathon_details"`
		Extracurriculars       *string `json:"extra_curriculars"` // Frontend uses extra_curriculars
		ClubMemberships        *string `json:"club_memberships"`
		Sports                 *string `json:"sports"`
		VolunteerWork          *string `json:"volunteer_work"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	tx, err := app.models.DB.Begin()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer tx.Rollback()

	achievements := &models.StudentAchievements{
		StudentID:              claims.UserID,
		Certifications:         input.Certifications,
		Awards:                 input.Awards,
		Workshops:              input.Workshops,
		Internships:            input.Internships,
		Projects:               input.Projects,
		LeetcodeProfile:        input.LeetcodeProfile,
		HackerrankProfile:      input.HackerrankProfile,
		CodeforcesProfile:      input.CodeforcesProfile,
		CodechefProfile:        input.CodechefProfile,
		LeetcodeRating:         input.LeetcodeRating,
		ProblemsSolved:         input.ProblemsSolved,
		HackathonsParticipated: input.HackathonsParticipated,
		HackathonsWon:          input.HackathonsWon,
		HackathonDetails:       input.HackathonDetails,
		Extracurriculars:       input.Extracurriculars,
		ClubMemberships:        input.ClubMemberships,
		Sports:                 input.Sports,
		VolunteerWork:          input.VolunteerWork,
	}

	if err := app.models.Students.UpsertAchievements(tx, achievements); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := tx.Commit(); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "Achievements saved"}, nil)
}

// updateAspirations updates student aspirations
func (app *application) updateAspirations(w http.ResponseWriter, r *http.Request) {
	claims, err := app.authenticateStudent(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	// Custom input struct with frontend field aliases
	var input struct {
		DreamCompanies     *string `json:"dream_company"` // Frontend uses dream_company (singular)
		PreferredRoles     *string `json:"preferred_roles"`
		PreferredLocations *string `json:"preferred_locations"`
		ExpectedSalary     *string `json:"expected_package"` // Frontend uses expected_package
		WillingToRelocate  bool    `json:"willing_to_relocate"`
		CareerObjective    *string `json:"career_goals"` // Frontend uses career_goals
		ShortTermGoals     *string `json:"short_term_goals"`
		LongTermGoals      *string `json:"long_term_goals"`
		HigherStudies      *string `json:"higher_studies"` // Extra field from frontend
		Strengths          *string `json:"strengths"`
		Weaknesses         *string `json:"weaknesses"`
		Hobbies            *string `json:"hobbies"`
		LanguagesKnown     *string `json:"languages_known"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	tx, err := app.models.DB.Begin()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer tx.Rollback()

	// Convert expected_package (number) to string for expected_salary
	aspirations := &models.StudentAspirations{
		StudentID:          claims.UserID,
		DreamCompanies:     input.DreamCompanies,
		PreferredRoles:     input.PreferredRoles,
		PreferredLocations: input.PreferredLocations,
		ExpectedSalary:     input.ExpectedSalary,
		WillingToRelocate:  input.WillingToRelocate,
		CareerObjective:    input.CareerObjective,
		ShortTermGoals:     input.ShortTermGoals,
		LongTermGoals:      input.LongTermGoals,
		Strengths:          input.Strengths,
		Weaknesses:         input.Weaknesses,
		Hobbies:            input.Hobbies,
		LanguagesKnown:     input.LanguagesKnown,
	}

	if err := app.models.Students.UpsertAspirations(tx, aspirations); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := tx.Commit(); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "Aspirations saved"}, nil)
}

// updateSkills updates student skills
func (app *application) updateSkills(w http.ResponseWriter, r *http.Request) {
	claims, err := app.authenticateStudent(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	var input struct {
		Skills []models.StudentSkill `json:"skills"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	tx, err := app.models.DB.Begin()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer tx.Rollback()

	if err := app.models.Students.UpsertSkills(tx, claims.UserID, input.Skills); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := tx.Commit(); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "Skills saved"}, nil)
}

// completeProfile marks the profile as complete
func (app *application) completeProfile(w http.ResponseWriter, r *http.Request) {
	claims, err := app.authenticateStudent(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	tx, err := app.models.DB.Begin()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer tx.Rollback()

	if err := app.models.Students.SetProfileCompleted(tx, claims.UserID); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := tx.Commit(); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "Profile marked as complete"}, nil)
}

// uploadPhoto handles profile photo upload
func (app *application) uploadPhoto(w http.ResponseWriter, r *http.Request) {
	claims, err := app.authenticateStudent(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	// Parse multipart form (max 5MB)
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		// Check if it's base64 JSON upload
		var input struct {
			Photo string `json:"photo"` // base64 encoded image
		}
		if err := app.readJSON(w, r, &input); err != nil {
			app.badRequestResponse(w, r, fmt.Errorf("unable to parse upload: %v", err))
			return
		}

		// Handle base64 upload
		if input.Photo != "" {
			// Validate and save base64 image
			photoURL, err := app.saveBase64Photo(claims.UserID, input.Photo)
			if err != nil {
				app.badRequestResponse(w, r, err)
				return
			}

			// Update student photo URL
			student, err := app.models.Students.GetByID(claims.UserID)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
			student.PhotoURL = &photoURL
			if err := app.models.Students.UpdateBasicInfo(student); err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			app.writeJSON(w, http.StatusOK, envelope{
				"message":   "Photo uploaded successfully",
				"photo_url": photoURL,
			}, nil)
			return
		}

		app.badRequestResponse(w, r, errors.New("no photo provided"))
		return
	}

	// Handle multipart form upload
	file, header, err := r.FormFile("photo")
	if err != nil {
		app.badRequestResponse(w, r, errors.New("photo file is required"))
		return
	}
	defer file.Close()

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		app.badRequestResponse(w, r, errors.New("file must be an image"))
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	filename := fmt.Sprintf("%d_%d%s", claims.UserID, time.Now().Unix(), ext)

	// Create uploads directory if not exists
	uploadDir := "./uploads/photos"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Save file
	filePath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Update student photo URL
	photoURL := "/uploads/photos/" + filename
	student, err := app.models.Students.GetByID(claims.UserID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	student.PhotoURL = &photoURL

	if err := app.models.Students.UpdateBasicInfo(student); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{
		"message":   "Photo uploaded successfully",
		"photo_url": photoURL,
	}, nil)
}

// saveBase64Photo saves a base64 encoded photo and returns the URL
func (app *application) saveBase64Photo(studentID int64, base64Data string) (string, error) {
	// Remove data URL prefix if present
	parts := strings.Split(base64Data, ",")
	var data string
	var ext string = ".jpg"

	if len(parts) == 2 {
		// Has prefix like "data:image/png;base64,"
		if strings.Contains(parts[0], "png") {
			ext = ".png"
		} else if strings.Contains(parts[0], "gif") {
			ext = ".gif"
		} else if strings.Contains(parts[0], "webp") {
			ext = ".webp"
		}
		data = parts[1]
	} else {
		data = base64Data
	}

	// Decode base64
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("invalid base64 data: %v", err)
	}

	// Create uploads directory
	uploadDir := "./uploads/photos"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}

	// Generate filename
	filename := fmt.Sprintf("%d_%d%s", studentID, time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// Write file
	if err := os.WriteFile(filePath, decoded, 0644); err != nil {
		return "", err
	}

	return "/uploads/photos/" + filename, nil
}

// ============================================
// COMMON ROUTES
// ============================================

// getSkills returns all available skills
func (app *application) getSkills(w http.ResponseWriter, r *http.Request) {
	skills, err := app.models.Skills.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Group by category
	grouped, err := app.models.Skills.GetGroupedByCategory()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{
		"skills":  skills,
		"grouped": grouped,
	}, nil)
}

// getBatches returns all batches
func (app *application) getBatches(w http.ResponseWriter, r *http.Request) {
	batches, err := app.models.Analytics.GetBatches()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"batches": batches}, nil)
}
