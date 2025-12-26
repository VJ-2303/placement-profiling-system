package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/VJ-2303/placement-profiling-system/internal/models"
)

// loginHandler initiates Microsoft OAuth flow
func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	// Generate random state for CSRF protection
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	state := base64.URLEncoding.EncodeToString(b)

	// Store state in cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		MaxAge:   300, // 5 minutes
		HttpOnly: true,
		Secure:   app.config.env == "production",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	// Redirect to Microsoft OAuth
	url := app.msOAuth.GetAuthURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// callbackHandler handles Microsoft OAuth callback
func (app *application) callbackHandler(w http.ResponseWriter, r *http.Request) {
	// Verify state
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		app.errorRedirect(w, r, "Session expired. Please try logging in again.")
		return
	}

	state := r.URL.Query().Get("state")
	if state != stateCookie.Value {
		app.errorRedirect(w, r, "Invalid session. Please try logging in again.")
		return
	}

	// Clear state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
	})

	// Check for OAuth errors
	if errMsg := r.URL.Query().Get("error"); errMsg != "" {
		errDesc := r.URL.Query().Get("error_description")
		app.logger.Printf("OAuth error: %s - %s", errMsg, errDesc)
		app.errorRedirect(w, r, "Authentication failed. Please try again.")
		return
	}

	// Exchange code for token
	code := r.URL.Query().Get("code")
	if code == "" {
		app.errorRedirect(w, r, "Authorization code not found.")
		return
	}

	token, err := app.msOAuth.Exchange(r.Context(), code)
	if err != nil {
		app.logger.Printf("Token exchange error: %v", err)
		app.errorRedirect(w, r, "Failed to authenticate. Please try again.")
		return
	}

	// Get user info from Microsoft
	userInfo, err := app.msOAuth.GetUserInfo(r.Context(), token)
	if err != nil {
		app.logger.Printf("User info error: %v", err)
		app.errorRedirect(w, r, "Failed to get user information.")
		return
	}

	app.logger.Printf("OAuth callback for user: %s (%s)", userInfo.DisplayName, userInfo.Mail)

	// Get email
	email := userInfo.Mail
	if email == "" {
		email = userInfo.UserPrincipalName
	}

	// Validate email domain for students
	emailLower := strings.ToLower(email)
	isAllowedDomain := strings.HasSuffix(emailLower, "@"+app.config.allowedDomain)

	// Check if user is an admin (pre-registered)
	admin, err := app.models.Admins.GetByEmail(email)
	if err == nil && admin != nil {
		// User is an admin
		app.logger.Printf("Admin login: %s", admin.Email)

		jwtToken, err := app.jwtService.GenerateToken(admin.ID, admin.Email, "admin")
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		redirectURL := fmt.Sprintf("%s/auth/callback.html?token=%s&role=admin", app.config.frontend.url, jwtToken)
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
		return
	}

	// Not an admin - must be a student
	if !isAllowedDomain {
		app.logger.Printf("Unauthorized domain: %s", email)
		app.errorRedirect(w, r, fmt.Sprintf("Only @%s email addresses are allowed for students.", app.config.allowedDomain))
		return
	}

	// Check if student exists
	student, err := app.models.Students.GetByEmail(email)
	if err != nil {
		if !errors.Is(err, models.ErrRecordNotFound) {
			app.serverErrorResponse(w, r, err)
			return
		}

		// Create new student
		student = &models.Student{
			Name:          userInfo.DisplayName,
			OfficialEmail: email,
		}

		err = app.models.Students.Insert(student)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		app.logger.Printf("New student created: %s", email)
	}

	// Update last login
	_ = app.models.Students.UpdateLastLogin(student.ID)

	// Generate JWT token
	jwtToken, err := app.jwtService.GenerateToken(student.ID, student.OfficialEmail, "student")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.Printf("Student login: %s", email)

	redirectURL := fmt.Sprintf("%s/auth/callback.html?token=%s&role=student", app.config.frontend.url, jwtToken)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

// getCurrentUser returns the current user's info from JWT
func (app *application) getCurrentUser(w http.ResponseWriter, r *http.Request) {
	claims, err := app.extractAndValidateToken(r)
	if err != nil {
		app.unauthorizedResponse(w, r)
		return
	}

	if claims.Role == "admin" {
		admin, err := app.models.Admins.GetByID(claims.UserID)
		if err != nil {
			if errors.Is(err, models.ErrRecordNotFound) {
				app.unauthorizedResponse(w, r)
				return
			}
			app.serverErrorResponse(w, r, err)
			return
		}

		app.writeJSON(w, http.StatusOK, envelope{
			"user": admin,
			"role": "admin",
		}, nil)
		return
	}

	// Student
	student, err := app.models.Students.GetByID(claims.UserID)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			app.unauthorizedResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{
		"user": student,
		"role": "student",
	}, nil)
}

// errorRedirect redirects to frontend with error message
func (app *application) errorRedirect(w http.ResponseWriter, r *http.Request, message string) {
	redirectURL := fmt.Sprintf("%s/auth/callback.html?error=%s", app.config.frontend.url, message)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}
