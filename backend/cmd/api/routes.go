package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/auth/login", app.loginHandler)
	router.HandlerFunc(http.MethodGet, "/auth/callback", app.callbackHandler)
	router.HandlerFunc(http.MethodGet, "/profile", app.profileHandler)

	// New comprehensive profile endpoints
	router.HandlerFunc(http.MethodPost, "/profile/complete", app.createStudentProfileHandler)
	router.HandlerFunc(http.MethodGet, "/profile/complete", app.getStudentProfileHandler)
	router.HandlerFunc(http.MethodPut, "/profile/complete", app.updateStudentProfileHandler)

	return app.enableCORS(router)
}
