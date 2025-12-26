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

	var input models.StudentPersonalDetails
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	input.StudentID = claims.UserID

	tx, err := app.models.DB.Begin()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer tx.Rollback()

	if err := app.models.Students.UpsertPersonalDetails(tx, &input); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := tx.Commit(); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "Personal details saved", "data": input}, nil)
}

// updateFamilyDetails updates student family details
func (app *application) updateFamilyDetails(w http.ResponseWriter, r *http.Request) {
	claims, err := app.authenticateStudent(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	var input models.StudentFamilyDetails
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	input.StudentID = claims.UserID

	tx, err := app.models.DB.Begin()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer tx.Rollback()

	if err := app.models.Students.UpsertFamilyDetails(tx, &input); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := tx.Commit(); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "Family details saved", "data": input}, nil)
}

// updateAcademics updates student academic details
func (app *application) updateAcademics(w http.ResponseWriter, r *http.Request) {
	claims, err := app.authenticateStudent(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	var input models.StudentAcademics
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	input.StudentID = claims.UserID

	tx, err := app.models.DB.Begin()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer tx.Rollback()

	if err := app.models.Students.UpsertAcademics(tx, &input); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := tx.Commit(); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "Academic details saved", "data": input}, nil)
}

// updateAchievements updates student achievements
func (app *application) updateAchievements(w http.ResponseWriter, r *http.Request) {
	claims, err := app.authenticateStudent(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	var input models.StudentAchievements
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	input.StudentID = claims.UserID

	tx, err := app.models.DB.Begin()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer tx.Rollback()

	if err := app.models.Students.UpsertAchievements(tx, &input); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := tx.Commit(); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "Achievements saved", "data": input}, nil)
}

// updateAspirations updates student aspirations
func (app *application) updateAspirations(w http.ResponseWriter, r *http.Request) {
	claims, err := app.authenticateStudent(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	var input models.StudentAspirations
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	input.StudentID = claims.UserID

	tx, err := app.models.DB.Begin()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer tx.Rollback()

	if err := app.models.Students.UpsertAspirations(tx, &input); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := tx.Commit(); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "Aspirations saved", "data": input}, nil)
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
