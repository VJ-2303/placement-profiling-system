package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/auth/login", app.loginHandler)
	router.HandlerFunc(http.MethodGet, "/auth/callback", app.callbackHandler)
	router.HandlerFunc(http.MethodGet, "/profile", app.StudentprofileHandler)

	router.HandlerFunc(http.MethodGet, "/admin/profile", app.AdminProfileHandler)
	router.HandlerFunc(http.MethodGet, "/admin/student/rollno/:rollno", app.SearchByRollNo)

	router.HandlerFunc(http.MethodPost, "/profile/complete", app.createStudentProfileHandler)
	router.HandlerFunc(http.MethodGet, "/profile/complete", app.getStudentProfileHandler)
	router.HandlerFunc(http.MethodPut, "/profile/complete", app.createStudentProfileHandler)

	return app.enableCORS(router)
}
