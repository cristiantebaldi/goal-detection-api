package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cristiantebaldi/goal-detection-api/pkg/formatter"
)

func HandleJSONError(response http.ResponseWriter, err error) {
	var unmarshalTypeError *json.UnmarshalTypeError
	var syntaxError *json.SyntaxError

	switch {
	case errors.As(err, &unmarshalTypeError):
		msg := fmt.Sprintf("The field '%s' needs a value of type '%s', but received '%s'",
			unmarshalTypeError.Field, unmarshalTypeError.Type, unmarshalTypeError.Value)
		formatter.HttpErrorResponse(response, formatter.ErrUnprocessableEntity, msg)

	case errors.As(err, &syntaxError):
		msg := fmt.Sprintf("Bad syntax at offset %d", syntaxError.Offset)
		formatter.HttpErrorResponse(response, formatter.ErrBadRequest, msg)

	case strings.HasPrefix(err.Error(), "json: unknown field"):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		msg := fmt.Sprintf("The field '%s' is not allowed", fieldName)
		formatter.HttpErrorResponse(response, formatter.ErrBadRequest, msg)

	case errors.Is(err, io.EOF):
		formatter.HttpErrorResponse(response, formatter.ErrBadRequest, "The body is empty")

	default:
		formatter.HttpErrorResponse(response, formatter.ErrBadRequest, "Unprocessed JSON")
	}
}
