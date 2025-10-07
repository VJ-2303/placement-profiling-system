package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/VJ-2303/placement-profiling-system/internal/models"
)

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	// Generate a random state for CSRF protection
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	state := base64.URLEncoding.EncodeToString(b)

	// Store state in a secure HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		MaxAge:   300, // 5 minutes
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	// Get the authorization URL and redirect user to Microsoft
	url := app.msOAuth.GetAuthURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (app *application) callbackHandler(w http.ResponseWriter, r *http.Request) {
	// Verify state parameter
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		app.badRequestResponse(w, r, errors.New("state cookie not found"))
		return
	}

	state := r.URL.Query().Get("state")
	if state != stateCookie.Value {
		app.badRequestResponse(w, r, errors.New("invalid state parameter"))
		return
	}

	// Clear the state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
	})

	// Get the authorization code
	code := r.URL.Query().Get("code")
	if code == "" {
		app.badRequestResponse(w, r, errors.New("authorization code not found"))
		return
	}

	// Exchange the authorization code for a token
	token, err := app.msOAuth.Exchange(r.Context(), code)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	userInfo, err := app.msOAuth.GetUserInfo(r.Context(), token)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.Printf("User Info: %+v\n", userInfo)

	// Determine the email to use
	email := userInfo.Mail
	if email == "" {
		email = userInfo.UserPrincipalName
	}

	// Check if student exists in database
	student, err := app.models.Students.GetByEmail(email)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			// Create new student
			newStudent := &models.Student{
				Name:          userInfo.DisplayName,
				OfficialEmail: email,
			}

			err = app.models.Students.Insert(newStudent)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			student = newStudent

		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	// Generate JWT token
	jwtToken, err := app.jwtService.GenerateToken(student.ID, student.OfficialEmail, "student")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	frontendURL := fmt.Sprintf("%s?token=%s", app.config.frontend.successURL, jwtToken)
	http.Redirect(w, r, frontendURL, http.StatusTemporaryRedirect)
}

func (app *application) profileHandler(w http.ResponseWriter, r *http.Request) {
	// Authenticate the request
	claims, err := app.authenticateRequest(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	// Get student from database
	student, err := app.models.Students.GetByID(claims.StudentID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Return student profile
	err = app.writeJSON(w, http.StatusOK, envelope{"student": student}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
