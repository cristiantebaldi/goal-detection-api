package util

import (
	"errors"
	"net/http"

	"github.com/cristiantebaldi/goal-detection-api/pkg/formatter"
)

type ConflictError interface {
	error
	IsConflict()
}

type NotFoundError interface {
	error
	IsNotFound()
}

type UnauthorizedError interface {
	error
	IsUnauthorized()
}

type ForbiddenError interface {
	error
	IsForbidden()
}

type BadRequestError interface {
	error
	IsBadRequest()
}

type InternalServerError interface {
	error
	IsInternalServerError()
}

func HandleAPIError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	if _, ok := errors.AsType[UnauthorizedError](err); ok {
		formatter.HttpErrorResponse(w, formatter.ErrAuth, err.Error())
		return
	}
	if _, ok := errors.AsType[ForbiddenError](err); ok {
		formatter.HttpErrorResponse(w, formatter.ErrForbidden, err.Error())
		return
	}
	if _, ok := errors.AsType[NotFoundError](err); ok {
		formatter.HttpErrorResponse(w, formatter.ErrNotFound, err.Error())
		return
	}
	if _, ok := errors.AsType[ConflictError](err); ok {
		formatter.HttpErrorResponse(w, formatter.ErrConflict, err.Error())
		return
	}
	if _, ok := errors.AsType[BadRequestError](err); ok {
		formatter.HttpErrorResponse(w, formatter.ErrBadRequest, err.Error())
		return
	}
	if _, ok := errors.AsType[InternalServerError](err); ok {
		formatter.HttpErrorResponse(w, formatter.ErrInternalServer, err.Error())
		return
	}

	formatter.HttpErrorResponse(w, formatter.ErrInternalServer, err.Error())
}
