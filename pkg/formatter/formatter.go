package formatter

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type errorAPIError struct {
	status int
	err    error
}

func (e errorAPIError) Error() string {
	return e.err.Error()
}

type HttpResponse struct {
	Message string `json:"message"`
}

func (e errorAPIError) HttpErrorResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.status)
	json.NewEncoder(w).Encode(HttpResponse{
		Message: e.err.Error(),
	})
}

func HttpErrorResponse(w http.ResponseWriter, err error, messages ...string) {
	if err == nil {
		return
	}

	var errorMessage strings.Builder
	if len(messages) > 0 {
		for i, msg := range messages {
			errorMessage.WriteString(msg)
			if i != len(messages)-1 {
				errorMessage.WriteString("\n")
			}
		}
	}

	if !isErrorAPIError(err) {
		if errorMessage.String() != "" {
			err = &errorAPIError{status: http.StatusInternalServerError, err: errors.New(errorMessage.String())}
		} else {
			err = &errorAPIError{status: http.StatusInternalServerError, err: err}
		}
	}

	if errorMessage.String() != "" {
		err = &errorAPIError{status: err.(*errorAPIError).status, err: errors.New(errorMessage.String())}
	}

	err.(*errorAPIError).HttpErrorResponse(w)
}

var (
	ErrAuth                = &errorAPIError{status: http.StatusUnauthorized, err: errors.New("not authorized")}
	ErrForbidden           = &errorAPIError{status: http.StatusForbidden, err: errors.New("forbidden")}
	ErrNotFound            = &errorAPIError{status: http.StatusNotFound, err: errors.New("resource not found")}
	ErrInternalServer      = &errorAPIError{status: http.StatusInternalServerError, err: errors.New("server error")}
	ErrBadRequest          = &errorAPIError{status: http.StatusBadRequest, err: errors.New("fill the fields")}
	ErrUnprocessableEntity = &errorAPIError{status: http.StatusUnprocessableEntity, err: errors.New("unprocessable entity")}
	ErrConflict            = &errorAPIError{status: http.StatusConflict, err: errors.New("conflict error")}
)

func isErrorAPIError(err error) bool {
	_, ok := err.(*errorAPIError)
	return ok
}
