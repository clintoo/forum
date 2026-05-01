package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"forum/backend/internal/dto"
	"forum/backend/internal/services"
)

// APIFunc is a handler signature that returns an error instead of writing it directly.
type APIFunc func(w http.ResponseWriter, r *http.Request) error

// ErrorHandlerAdapter converts an APIFunc into a standard http.HandlerFunc.
// It executes the handler and, on error, logs the internal detail and writes a
// client-safe error response derived from getAPIError.
func ErrorHandlerAdapter(h APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			status, msg := getAPIError(err)
			log.Printf("[ERROR] %v", err)
			writeJSON(w, status, dto.Response{Status: "error", Message: msg})
		}
	}
}

// getAPIError maps an error to a safe HTTP status code and a client-safe message.
// Internal/unknown errors are masked to prevent leaking implementation details.
func getAPIError(err error) (int, string) {
	var validationErr *services.ValidationError
	var conflictErr *services.ConflictError
	var notFoundErr *services.NotFoundError
	var internalErr *services.InternalError

	switch {
	case errors.As(err, &validationErr):
		return http.StatusBadRequest, validationErr.Msg
	case errors.As(err, &conflictErr):
		return http.StatusConflict, conflictErr.Msg
	case errors.As(err, &notFoundErr):
		return http.StatusNotFound, notFoundErr.Msg
	case errors.As(err, &internalErr):
		return http.StatusInternalServerError, "An unexpected error occurred"
	default:
		return http.StatusInternalServerError, "An unexpected error occurred"
	}
}

// writeJSON encodes v as JSON and writes it with the given HTTP status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// writeError writes a JSON error response using a standardized status DTO.
func writeError(w http.ResponseWriter, message string, status int) {
	writeJSON(w, status, dto.Response{Status: "error", Message: message})
}

// writeSuccess writes a standardized success response with optional data.
func writeSuccess(w http.ResponseWriter, status int, message string, data any) {
	writeJSON(w, status, dto.Response{Status: "success", Message: message, Data: data})
}

// statusCodeFromServiceError returns the appropriate HTTP status for a service error.
func statusCodeFromServiceError(err error) int {
	status, _ := getAPIError(err)
	return status
}

// getURLParamInt safely retrieves an integer URL parameter and returns a validation error if missing or invalid.
func getURLParamInt(r *http.Request, param string) (int, error) {
	valStr := r.PathValue(param)
	if valStr == "" {
		return 0, services.NewValidationError(fmt.Sprintf("%s is required", param))
	}

	val, err := strconv.Atoi(valStr)
	if err != nil || val <= 0 {
		return 0, services.NewValidationError(fmt.Sprintf("invalid %s", param))
	}

	return val, nil
}

// safeTemplateExecute executes a template with buffering to prevent partial writes on error.
// This prevents the "superfluous response.WriteHeader call" error when template execution fails.
func safeTemplateExecute(w http.ResponseWriter, t *template.Template, name string, data any) error {
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, name, data); err != nil {
		return err
	}
	// If template executed successfully, write to response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := buf.WriteTo(w)
	return err
}
