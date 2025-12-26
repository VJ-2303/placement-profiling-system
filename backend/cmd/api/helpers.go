package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/VJ-2303/placement-profiling-system/internal/auth"
	"github.com/julienschmidt/httprouter"
)

type envelope map[string]interface{}

// ============================================
// JSON HELPERS
// ============================================

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 2_097_152 // 2MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON object")
	}

	return nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

// ============================================
// ERROR RESPONSES
// ============================================

func (app *application) errorResponse(w http.ResponseWriter, _ *http.Request, status int, message interface{}) {
	env := envelope{"error": message}
	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logger.Print(err)
		w.WriteHeader(500)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Printf("ERROR: %v", err)
	message := "The server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) unauthorizedResponse(w http.ResponseWriter, r *http.Request) {
	message := "You must be authenticated to access this resource"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	message := "You don't have permission to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}

func (app *application) validationErrorResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

// ============================================
// AUTHENTICATION HELPERS
// ============================================

func (app *application) extractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("authorization header format must be 'Bearer {token}'")
	}

	return parts[1], nil
}

func (app *application) extractAndValidateToken(r *http.Request) (*auth.Claims, error) {
	tokenString, err := app.extractTokenFromHeader(r)
	if err != nil {
		return nil, err
	}

	claims, err := app.jwtService.ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}

func (app *application) authenticateStudent(r *http.Request) (*auth.Claims, error) {
	claims, err := app.extractAndValidateToken(r)
	if err != nil {
		return nil, err
	}

	if claims.Role != "student" {
		return nil, errors.New("student access required")
	}

	return claims, nil
}

func (app *application) authenticateAdmin(r *http.Request) (*auth.Claims, error) {
	claims, err := app.extractAndValidateToken(r)
	if err != nil {
		return nil, err
	}

	if claims.Role != "admin" {
		return nil, errors.New("admin access required")
	}

	return claims, nil
}

// ============================================
// URL PARAMETER HELPERS
// ============================================

func (app *application) readIDParam(r *http.Request, paramName string) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	idStr := params.ByName(paramName)
	if idStr == "" {
		return 0, fmt.Errorf("missing %s parameter", paramName)
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id < 1 {
		return 0, fmt.Errorf("invalid %s parameter", paramName)
	}

	return id, nil
}

func (app *application) readStringParam(r *http.Request, paramName string) string {
	params := httprouter.ParamsFromContext(r.Context())
	return params.ByName(paramName)
}

// ============================================
// QUERY STRING HELPERS
// ============================================

func (app *application) readString(qs map[string][]string, key string, defaultValue string) string {
	s := qs[key]
	if len(s) == 0 {
		return defaultValue
	}
	return s[0]
}

func (app *application) readInt(qs map[string][]string, key string, defaultValue int) int {
	s := qs[key]
	if len(s) == 0 {
		return defaultValue
	}
	i, err := strconv.Atoi(s[0])
	if err != nil {
		return defaultValue
	}
	return i
}

func (app *application) readFloat(qs map[string][]string, key string, defaultValue float64) float64 {
	s := qs[key]
	if len(s) == 0 {
		return defaultValue
	}
	f, err := strconv.ParseFloat(s[0], 64)
	if err != nil {
		return defaultValue
	}
	return f
}

func (app *application) readBool(qs map[string][]string, key string) *bool {
	s := qs[key]
	if len(s) == 0 {
		return nil
	}
	b, err := strconv.ParseBool(s[0])
	if err != nil {
		return nil
	}
	return &b
}
