package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// Health check
	router.HandlerFunc(http.MethodGet, "/health", app.healthCheckHandler)

	// ============================================
	// AUTH ROUTES
	// ============================================
	router.HandlerFunc(http.MethodGet, "/auth/login", app.loginHandler)
	router.HandlerFunc(http.MethodGet, "/auth/callback", app.callbackHandler)
	router.HandlerFunc(http.MethodGet, "/auth/me", app.getCurrentUser)

	// ============================================
	// STUDENT ROUTES
	// ============================================
	router.HandlerFunc(http.MethodGet, "/api/student/profile", app.getStudentProfile)
	router.HandlerFunc(http.MethodPut, "/api/student/profile", app.updateStudentProfile)
	router.HandlerFunc(http.MethodPost, "/api/student/profile/personal", app.updatePersonalDetails)
	router.HandlerFunc(http.MethodPost, "/api/student/profile/family", app.updateFamilyDetails)
	router.HandlerFunc(http.MethodPost, "/api/student/profile/academics", app.updateAcademics)
	router.HandlerFunc(http.MethodPost, "/api/student/profile/achievements", app.updateAchievements)
	router.HandlerFunc(http.MethodPost, "/api/student/profile/aspirations", app.updateAspirations)
	router.HandlerFunc(http.MethodPost, "/api/student/profile/skills", app.updateSkills)
	router.HandlerFunc(http.MethodPost, "/api/student/profile/complete", app.completeProfile)
	router.HandlerFunc(http.MethodPost, "/api/student/photo", app.uploadPhoto)

	// ============================================
	// ADMIN ROUTES
	// ============================================

	// Dashboard & Analytics
	router.HandlerFunc(http.MethodGet, "/api/admin/dashboard", app.getDashboard)
	router.HandlerFunc(http.MethodGet, "/api/admin/analytics/batch", app.getBatchStats)
	router.HandlerFunc(http.MethodGet, "/api/admin/analytics/skills", app.getSkillStats)
	router.HandlerFunc(http.MethodGet, "/api/admin/analytics/cgpa", app.getCGPADistribution)
	router.HandlerFunc(http.MethodGet, "/api/admin/analytics/companies", app.getCompanyStats)
	router.HandlerFunc(http.MethodGet, "/api/admin/activity", app.getRecentActivity)

	// Student Management
	router.HandlerFunc(http.MethodGet, "/api/admin/students", app.listStudents)
	router.HandlerFunc(http.MethodGet, "/api/admin/students/export", app.exportStudentsCSV)
	router.HandlerFunc(http.MethodGet, "/api/admin/students/:id", app.getStudentByID)
	router.HandlerFunc(http.MethodGet, "/api/admin/students/roll/:rollno", app.getStudentByRollNo)
	router.HandlerFunc(http.MethodPut, "/api/admin/students/:id/status", app.updateStudentStatus)
	router.HandlerFunc(http.MethodPatch, "/api/admin/students/:id/status", app.updateStudentStatus)

	// Placement Management
	router.HandlerFunc(http.MethodGet, "/api/admin/placements", app.listPlacements)
	router.HandlerFunc(http.MethodPost, "/api/admin/placements", app.createPlacement)
	router.HandlerFunc(http.MethodPut, "/api/admin/placements/:id", app.updatePlacement)
	router.HandlerFunc(http.MethodDelete, "/api/admin/placements/:id", app.deletePlacement)

	// Company Management
	router.HandlerFunc(http.MethodGet, "/api/admin/companies", app.listCompanies)
	router.HandlerFunc(http.MethodPost, "/api/admin/companies", app.createCompany)
	router.HandlerFunc(http.MethodPut, "/api/admin/companies/:id", app.updateCompany)
	router.HandlerFunc(http.MethodDelete, "/api/admin/companies/:id", app.deleteCompany)
	router.HandlerFunc(http.MethodGet, "/api/admin/companies/search", app.searchCompanies)

	// ============================================
	// COMMON ROUTES
	// ============================================
	router.HandlerFunc(http.MethodGet, "/api/skills", app.getSkills)
	router.HandlerFunc(http.MethodGet, "/api/batches", app.getBatches)

	return app.enableCORS(router)
}

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	app.writeJSON(w, http.StatusOK, envelope{
		"status": "ok",
		"env":    app.config.env,
	}, nil)
}
