package main

import (
	"errors"
	"net/http"

	"github.com/VJ-2303/placement-profiling-system/internal/models"
	"github.com/julienschmidt/httprouter"
)

func (app *application) SearchByRollNo(w http.ResponseWriter, r *http.Request) {

	_, err := app.authenticateAdmin(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	params := httprouter.ParamsFromContext(r.Context())

	rollNo := params.ByName("rollno")

	studentProfile, err := app.models.Admins.GetFullProfileByRollNo(rollNo)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"student": studentProfile}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
