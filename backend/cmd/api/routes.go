package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter()

	// Health check
	router.HandleFunc("/health", app.healthCheckHandler).Methods(http.MethodGet)

	// ============================================
	// AUTH ROUTES
	// ============================================
	router.HandleFunc("/auth/login", app.loginHandler).Methods(http.MethodGet)
	router.HandleFunc("/auth/callback", app.callbackHandler).Methods(http.MethodGet)
	router.HandleFunc("/auth/me", app.getCurrentUser).Methods(http.MethodGet)

	// ============================================
	// STUDENT ROUTES
	// ============================================
	router.HandleFunc("/api/student/profile", app.getStudentProfile).Methods(http.MethodGet)
	router.HandleFunc("/api/student/profile", app.updateStudentProfile).Methods(http.MethodPut)
	router.HandleFunc("/api/student/profile/personal", app.updatePersonalDetails).Methods(http.MethodPut)
	router.HandleFunc("/api/student/profile/family", app.updateFamilyDetails).Methods(http.MethodPut)
	router.HandleFunc("/api/student/profile/academics", app.updateAcademics).Methods(http.MethodPut)
	router.HandleFunc("/api/student/profile/achievements", app.updateAchievements).Methods(http.MethodPut)
	router.HandleFunc("/api/student/profile/aspirations", app.updateAspirations).Methods(http.MethodPut)
	router.HandleFunc("/api/student/profile/skills", app.updateSkills).Methods(http.MethodPut)
	router.HandleFunc("/api/student/profile/complete", app.completeProfile).Methods(http.MethodPost)
	router.HandleFunc("/api/student/photo", app.uploadPhoto).Methods(http.MethodPost)

	// ============================================
	// ADMIN ROUTES
	// ============================================

	// Dashboard & Analytics
	router.HandleFunc("/api/admin/dashboard", app.getDashboard).Methods(http.MethodGet)
	router.HandleFunc("/api/admin/analytics/batch", app.getBatchStats).Methods(http.MethodGet)
	router.HandleFunc("/api/admin/analytics/skills", app.getSkillStats).Methods(http.MethodGet)
	router.HandleFunc("/api/admin/analytics/cgpa", app.getCGPADistribution).Methods(http.MethodGet)
	router.HandleFunc("/api/admin/analytics/companies", app.getCompanyStats).Methods(http.MethodGet)
	router.HandleFunc("/api/admin/activity", app.getRecentActivity).Methods(http.MethodGet)

	// Student Management - IMPORTANT: specific routes before parameterized routes
	router.HandleFunc("/api/admin/students/export", app.exportStudentsCSV).Methods(http.MethodGet)
	router.HandleFunc("/api/admin/students/roll/{rollno}", app.getStudentByRollNo).Methods(http.MethodGet)
	router.HandleFunc("/api/admin/students/{id:[0-9]+}/status", app.updateStudentStatus).Methods(http.MethodPut, http.MethodPatch)
	router.HandleFunc("/api/admin/students/{id:[0-9]+}", app.getStudentByID).Methods(http.MethodGet)
	router.HandleFunc("/api/admin/students", app.listStudents).Methods(http.MethodGet)

	// Placement Management
	router.HandleFunc("/api/admin/placements/{id:[0-9]+}", app.updatePlacement).Methods(http.MethodPut)
	router.HandleFunc("/api/admin/placements/{id:[0-9]+}", app.deletePlacement).Methods(http.MethodDelete)
	router.HandleFunc("/api/admin/placements", app.listPlacements).Methods(http.MethodGet)
	router.HandleFunc("/api/admin/placements", app.createPlacement).Methods(http.MethodPost)

	// Company Management - specific routes before parameterized routes
	router.HandleFunc("/api/admin/companies/search", app.searchCompanies).Methods(http.MethodGet)
	router.HandleFunc("/api/admin/companies/{id:[0-9]+}", app.updateCompany).Methods(http.MethodPut)
	router.HandleFunc("/api/admin/companies/{id:[0-9]+}", app.deleteCompany).Methods(http.MethodDelete)
	router.HandleFunc("/api/admin/companies", app.listCompanies).Methods(http.MethodGet)
	router.HandleFunc("/api/admin/companies", app.createCompany).Methods(http.MethodPost)

	// ============================================
	// COMMON ROUTES
	// ============================================
	router.HandleFunc("/api/skills", app.getSkills).Methods(http.MethodGet)
	router.HandleFunc("/api/batches", app.getBatches).Methods(http.MethodGet)

	return app.enableCORS(router)
}

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	app.writeJSON(w, http.StatusOK, envelope{
		"status":  "healthy",
		"service": "placement-api",
		"env":     app.config.env,
	}, nil)
}
