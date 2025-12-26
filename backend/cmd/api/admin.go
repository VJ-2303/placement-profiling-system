package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/VJ-2303/placement-profiling-system/internal/models"
)

// ============================================
// DASHBOARD & ANALYTICS
// ============================================

// getDashboard returns admin dashboard data
func (app *application) getDashboard(w http.ResponseWriter, r *http.Request) {
	claims, err := app.authenticateAdmin(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	// Get admin info
	admin, err := app.models.Admins.GetByID(claims.UserID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Parse optional batch filter
	qs := r.URL.Query()
	var batchYear *int
	if batchStr := qs.Get("batch"); batchStr != "" {
		b, err := strconv.Atoi(batchStr)
		if err == nil {
			batchYear = &b
		}
	}

	// Get dashboard stats
	stats, err := app.models.Analytics.GetDashboardStats(batchYear)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Get recent activity
	activity, _ := app.models.Analytics.GetRecentActivity(10)

	// Get batches
	batches, _ := app.models.Analytics.GetBatches()

	app.writeJSON(w, http.StatusOK, envelope{
		"admin":    admin,
		"stats":    stats,
		"activity": activity,
		"batches":  batches,
	}, nil)
}

// getBatchStats returns batch-wise statistics
func (app *application) getBatchStats(w http.ResponseWriter, r *http.Request) {
	_, err := app.authenticateAdmin(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	stats, err := app.models.Analytics.GetBatchWiseStats()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"batch_stats": stats}, nil)
}

// getSkillStats returns skill distribution statistics
func (app *application) getSkillStats(w http.ResponseWriter, r *http.Request) {
	_, err := app.authenticateAdmin(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	stats, err := app.models.Analytics.GetSkillStats()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"skill_stats": stats}, nil)
}

// getCGPADistribution returns CGPA distribution
func (app *application) getCGPADistribution(w http.ResponseWriter, r *http.Request) {
	_, err := app.authenticateAdmin(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	qs := r.URL.Query()
	var batchYear *int
	if batchStr := qs.Get("batch"); batchStr != "" {
		b, err := strconv.Atoi(batchStr)
		if err == nil {
			batchYear = &b
		}
	}

	stats, err := app.models.Analytics.GetCGPADistribution(batchYear)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"cgpa_distribution": stats}, nil)
}

// getCompanyStats returns placement statistics by company
func (app *application) getCompanyStats(w http.ResponseWriter, r *http.Request) {
	_, err := app.authenticateAdmin(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	qs := r.URL.Query()
	var batchYear *int
	if batchStr := qs.Get("batch"); batchStr != "" {
		b, err := strconv.Atoi(batchStr)
		if err == nil {
			batchYear = &b
		}
	}

	stats, err := app.models.Analytics.GetCompanyStats(batchYear)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"company_stats": stats}, nil)
}

// getRecentActivity returns recent activities
func (app *application) getRecentActivity(w http.ResponseWriter, r *http.Request) {
	_, err := app.authenticateAdmin(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	limit := app.readInt(r.URL.Query(), "limit", 20)
	activities, err := app.models.Analytics.GetRecentActivity(limit)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"activities": activities}, nil)
}

// ============================================
// STUDENT MANAGEMENT
// ============================================

// listStudents returns paginated list of students with filters
func (app *application) listStudents(w http.ResponseWriter, r *http.Request) {
	_, err := app.authenticateAdmin(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	qs := r.URL.Query()

	filter := models.StudentFilter{
		Search:   app.readString(qs, "search", ""),
		Page:     app.readInt(qs, "page", 1),
		PageSize: app.readInt(qs, "page_size", 20),
	}

	// Batch filter
	if batchStr := qs.Get("batch"); batchStr != "" {
		b, err := strconv.Atoi(batchStr)
		if err == nil {
			filter.BatchYear = &b
		}
	}

	// Placement status filter
	if status := qs.Get("status"); status != "" {
		ps := models.PlacementStatus(status)
		filter.PlacementStatus = &ps
	}

	// CGPA filters
	if minCGPA := app.readFloat(qs, "min_cgpa", 0); minCGPA > 0 {
		filter.MinCGPA = &minCGPA
	}
	if maxCGPA := app.readFloat(qs, "max_cgpa", 0); maxCGPA > 0 {
		filter.MaxCGPA = &maxCGPA
	}

	// Backlog filter
	filter.HasBacklogs = app.readBool(qs, "has_backlogs")

	result, err := app.models.Students.List(filter)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{
		"students":    result.Students,
		"total":       result.Total,
		"page":        result.Page,
		"page_size":   result.PageSize,
		"total_pages": result.TotalPages,
	}, nil)
}

// getStudentByID returns full profile of a student
func (app *application) getStudentByID(w http.ResponseWriter, r *http.Request) {
	_, err := app.authenticateAdmin(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	id, err := app.readIDParam(r, "id")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	profile, err := app.models.Students.GetFullProfile(id)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	// Get placement info
	placement, _ := app.models.Placements.GetByStudentID(id)
	profile.Placement = placement

	app.writeJSON(w, http.StatusOK, envelope{"profile": profile}, nil)
}

// getStudentByRollNo returns student by roll number
func (app *application) getStudentByRollNo(w http.ResponseWriter, r *http.Request) {
	_, err := app.authenticateAdmin(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	rollNo := app.readStringParam(r, "rollno")
	if rollNo == "" {
		app.badRequestResponse(w, r, errors.New("roll number is required"))
		return
	}

	student, err := app.models.Students.GetByRollNo(rollNo)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	profile, err := app.models.Students.GetFullProfile(student.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Get placement info
	placement, _ := app.models.Placements.GetByStudentID(student.ID)
	profile.Placement = placement

	app.writeJSON(w, http.StatusOK, envelope{"profile": profile}, nil)
}

// updateStudentStatus updates student placement status
func (app *application) updateStudentStatus(w http.ResponseWriter, r *http.Request) {
	_, err := app.authenticateAdmin(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	id, err := app.readIDParam(r, "id")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input struct {
		Status                 string `json:"status"`
		PlacementStatus        string `json:"placement_status"`
		IsEligibleForPlacement *bool  `json:"is_eligible_for_placement"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Accept both "status" and "placement_status" fields
	statusStr := input.Status
	if statusStr == "" {
		statusStr = input.PlacementStatus
	}

	if statusStr != "" {
		status := models.PlacementStatus(statusStr)
		if err := app.models.Students.UpdatePlacementStatus(id, status); err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "Student status updated"}, nil)
}

// exportStudentsCSV exports student data as CSV
func (app *application) exportStudentsCSV(w http.ResponseWriter, r *http.Request) {
	_, err := app.authenticateAdmin(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	qs := r.URL.Query()
	filter := models.StudentFilter{
		Page:     1,
		PageSize: 10000, // Export all
	}

	// Apply filters
	if batchStr := qs.Get("batch"); batchStr != "" {
		b, err := strconv.Atoi(batchStr)
		if err == nil {
			filter.BatchYear = &b
		}
	}
	if status := qs.Get("status"); status != "" {
		ps := models.PlacementStatus(status)
		filter.PlacementStatus = &ps
	}

	result, err := app.models.Students.List(filter)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Set CSV headers
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=students_export.csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header row
	header := []string{
		"ID", "Name", "Email", "Roll No", "Batch", "Profile Completed",
		"Placement Status", "CGPA", "Mobile", "Placed Company", "Package (LPA)",
	}
	writer.Write(header)

	// Write data rows
	for _, s := range result.Students {
		row := []string{
			fmt.Sprintf("%d", s.ID),
			s.Name,
			s.OfficialEmail,
			ptrToString(s.RollNo),
			ptrIntToString(s.BatchYear),
			fmt.Sprintf("%t", s.IsProfileCompleted),
			string(s.PlacementStatus),
			ptrFloatToString(s.CGPAOverall),
			ptrToString(s.MobileNumber),
			ptrToString(s.PlacedCompany),
			ptrFloatToString(s.PackageLPA),
		}
		writer.Write(row)
	}
}

// ============================================
// PLACEMENT MANAGEMENT
// ============================================

// listPlacements returns all placement records
func (app *application) listPlacements(w http.ResponseWriter, r *http.Request) {
	_, err := app.authenticateAdmin(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	placements, err := app.models.Placements.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"placements": placements}, nil)
}

// createPlacement creates a new placement record
func (app *application) createPlacement(w http.ResponseWriter, r *http.Request) {
	claims, err := app.authenticateAdmin(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	var input struct {
		StudentID   int64    `json:"student_id"`
		CompanyID   *int64   `json:"company_id"`
		CompanyName string   `json:"company_name"`
		JobRole     *string  `json:"job_role"`
		PackageLPA  *float64 `json:"package_lpa"`
		PackageCTC  *string  `json:"package_ctc"`
		JoiningDate *string  `json:"joining_date"`
		OfferDate   *string  `json:"offer_date"`
		OfferType   *string  `json:"offer_type"`
		JobLocation *string  `json:"job_location"`
		Remarks     *string  `json:"remarks"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate student exists
	_, err = app.models.Students.GetByID(input.StudentID)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			app.badRequestResponse(w, r, errors.New("student not found"))
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	placement := &models.PlacementRecord{
		StudentID:   input.StudentID,
		CompanyID:   input.CompanyID,
		CompanyName: input.CompanyName,
		JobRole:     input.JobRole,
		PackageLPA:  input.PackageLPA,
		PackageCTC:  input.PackageCTC,
		JoiningDate: input.JoiningDate,
		OfferDate:   input.OfferDate,
		OfferType:   input.OfferType,
		JobLocation: input.JobLocation,
		Remarks:     input.Remarks,
		IsAccepted:  true,
		VerifiedBy:  &claims.UserID,
	}

	if err := app.models.Placements.Insert(placement); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Update student status
	if err := app.models.Students.UpdatePlacementStatus(input.StudentID, models.PlacementStatusPlaced); err != nil {
		app.logger.Printf("Warning: Failed to update student placement status: %v", err)
	}

	app.writeJSON(w, http.StatusCreated, envelope{"placement": placement}, nil)
}

// updatePlacement updates a placement record
func (app *application) updatePlacement(w http.ResponseWriter, r *http.Request) {
	_, err := app.authenticateAdmin(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	id, err := app.readIDParam(r, "id")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input struct {
		CompanyID   *int64   `json:"company_id"`
		CompanyName string   `json:"company_name"`
		JobRole     *string  `json:"job_role"`
		PackageLPA  *float64 `json:"package_lpa"`
		PackageCTC  *string  `json:"package_ctc"`
		JoiningDate *string  `json:"joining_date"`
		OfferDate   *string  `json:"offer_date"`
		OfferType   *string  `json:"offer_type"`
		JobLocation *string  `json:"job_location"`
		IsAccepted  bool     `json:"is_accepted"`
		Remarks     *string  `json:"remarks"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	placement := &models.PlacementRecord{
		ID:          id,
		CompanyID:   input.CompanyID,
		CompanyName: input.CompanyName,
		JobRole:     input.JobRole,
		PackageLPA:  input.PackageLPA,
		PackageCTC:  input.PackageCTC,
		JoiningDate: input.JoiningDate,
		OfferDate:   input.OfferDate,
		OfferType:   input.OfferType,
		JobLocation: input.JobLocation,
		IsAccepted:  input.IsAccepted,
		Remarks:     input.Remarks,
	}

	if err := app.models.Placements.Update(placement); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"placement": placement}, nil)
}

// deletePlacement deletes a placement record
func (app *application) deletePlacement(w http.ResponseWriter, r *http.Request) {
	_, err := app.authenticateAdmin(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	id, err := app.readIDParam(r, "id")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.models.Placements.Delete(id); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "Placement record deleted"}, nil)
}

// ============================================
// COMPANY MANAGEMENT
// ============================================

// listCompanies returns all companies
func (app *application) listCompanies(w http.ResponseWriter, r *http.Request) {
	_, err := app.authenticateAdmin(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	companies, err := app.models.Companies.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"companies": companies}, nil)
}

// createCompany creates a new company
func (app *application) createCompany(w http.ResponseWriter, r *http.Request) {
	_, err := app.authenticateAdmin(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	var input models.Company
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name == "" {
		app.badRequestResponse(w, r, errors.New("company name is required"))
		return
	}

	if err := app.models.Companies.Insert(&input); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"company": input}, nil)
}

// updateCompany updates a company
func (app *application) updateCompany(w http.ResponseWriter, r *http.Request) {
	_, err := app.authenticateAdmin(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	id, err := app.readIDParam(r, "id")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	company, err := app.models.Companies.GetByID(id)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	var input struct {
		Name         *string `json:"name"`
		Website      *string `json:"website"`
		Industry     *string `json:"industry"`
		CompanyType  *string `json:"company_type"`
		Description  *string `json:"description"`
		HRName       *string `json:"hr_name"`
		HREmail      *string `json:"hr_email"`
		HRPhone      *string `json:"hr_phone"`
		Headquarters *string `json:"headquarters"`
		IsActive     *bool   `json:"is_active"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		company.Name = *input.Name
	}
	if input.Website != nil {
		company.Website = input.Website
	}
	if input.Industry != nil {
		company.Industry = input.Industry
	}
	if input.CompanyType != nil {
		company.CompanyType = input.CompanyType
	}
	if input.Description != nil {
		company.Description = input.Description
	}
	if input.HRName != nil {
		company.HRName = input.HRName
	}
	if input.HREmail != nil {
		company.HREmail = input.HREmail
	}
	if input.HRPhone != nil {
		company.HRPhone = input.HRPhone
	}
	if input.Headquarters != nil {
		company.Headquarters = input.Headquarters
	}
	if input.IsActive != nil {
		company.IsActive = *input.IsActive
	}

	if err := app.models.Companies.Update(company); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"company": company}, nil)
}

// searchCompanies searches companies by name
func (app *application) searchCompanies(w http.ResponseWriter, r *http.Request) {
	_, err := app.authenticateAdmin(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		app.badRequestResponse(w, r, errors.New("search query is required"))
		return
	}

	companies, err := app.models.Companies.Search(query)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"companies": companies}, nil)
}

// deleteCompany deletes a company
func (app *application) deleteCompany(w http.ResponseWriter, r *http.Request) {
	_, err := app.authenticateAdmin(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	id, err := app.readIDParam(r, "id")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.Companies.Delete(id)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "company deleted"}, nil)
}

// ============================================
// HELPER FUNCTIONS
// ============================================

func ptrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ptrIntToString(i *int) string {
	if i == nil {
		return ""
	}
	return fmt.Sprintf("%d", *i)
}

func ptrFloatToString(f *float64) string {
	if f == nil {
		return ""
	}
	return fmt.Sprintf("%.2f", *f)
}
