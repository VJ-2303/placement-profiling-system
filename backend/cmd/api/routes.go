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

	return app.enableCORS(router)
}
