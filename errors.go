package wrap

import (
	"fmt"
	"net/http"

	"github.com/torniker/wrap/logger"
)

// Error checks error type and responses accordingly
func (ctx *Ctx) Error(err error) {
	switch err.(type) {
	case ErrorBadRequest:
		ctx.Response.SetStatus(http.StatusBadRequest)
		e := err.(ErrorBadRequest)
		ctx.Response.Write(e)
	case ErrorStatusUnauthorized:
		ctx.Response.SetStatus(http.StatusUnauthorized)
		e := err.(ErrorStatusUnauthorized)
		ctx.Response.Write(e)
	case ErrorStatusNotAllowed:
		ctx.Response.SetStatus(http.StatusMethodNotAllowed)
		e := err.(ErrorStatusNotAllowed)
		ctx.Response.Write(e)
	case ErrorStatusNotFound:
		ctx.Response.SetStatus(http.StatusNotFound)
		e := err.(ErrorStatusNotFound)
		ctx.Response.Write(e)
	case ErrorInternalServerError:
		ctx.Response.SetStatus(http.StatusInternalServerError)
		e := err.(ErrorInternalServerError)
		ctx.Response.Write(e)
	case ErrorUnprocessableEntity:
		ctx.Response.SetStatus(http.StatusUnprocessableEntity)
		e := err.(ErrorUnprocessableEntity)
		ctx.Response.Write(e)
	default:
		ctx.Response.SetStatus(http.StatusInternalServerError)
		ctx.Response.Write(ErrorInternalServerError{Message: err.Error()})
	}
}

// NotFound response 404
func (ctx *Ctx) NotFound() error {
	e := ErrorStatusNotFound{
		Message:  "not found",
		Internal: fmt.Sprintf("url: %v not found", ctx.Request.Path().URL().Path),
	}
	logger.Error(e.Internal)
	return e
}

// Unauthorized response 401
func (ctx *Ctx) Unauthorized() error {
	e := ErrorStatusUnauthorized{
		Message:  "unauthorized",
		Internal: fmt.Sprintf("user: %v is unauthorized to request: %v", ctx.User, ctx.Request.Path().URL().Path),
	}
	logger.Error(e.Internal)
	return e
}

// InternalError response 404
func (ctx *Ctx) InternalError(err error) error {
	e := ErrorInternalServerError{
		Message:  "internal server error",
		Internal: err.Error(),
	}
	logger.Error(e.Internal)
	return e
}

// BadRequest response 404
func (ctx *Ctx) BadRequest() error {
	e := ErrorInternalServerError{
		Message:  "bad request",
		Internal: fmt.Sprintf("bad request: %#v", ctx.Request),
	}
	logger.Error(e.Internal)
	return e
}

// NotAllowed response 404
func (ctx *Ctx) NotAllowed() error {
	e := ErrorInternalServerError{
		Message:  "not allowed",
		Internal: fmt.Sprintf("user: %v is not allowed to request: %v", ctx.User, ctx.Request.Path().URL().Path),
	}
	logger.Error(e.Internal)
	return e
}

// UnprocessableEntity respnse 422
func (ctx *Ctx) UnprocessableEntity(code int, message string) error {
	e := ErrorUnprocessableEntity{
		Code:     code,
		Message:  message,
		Internal: fmt.Sprintf("UnprocessableEntity code: %v, message: %v ", code, message),
	}
	logger.Error(e.Internal)
	return e
}

// ErrorBadRequest type for bad request
type ErrorBadRequest struct {
	Message  string `json:"message"`
	Internal string `json:"-"`
}

func (e ErrorBadRequest) Error() string {
	return e.Message
}

// ErrorStatusUnauthorized type for Unauthorized
type ErrorStatusUnauthorized struct {
	Message  string `json:"message"`
	Internal string `json:"-"`
}

func (e ErrorStatusUnauthorized) Error() string {
	return e.Message
}

// ErrorStatusNotAllowed type for not allowed
type ErrorStatusNotAllowed struct {
	Message  string `json:"message"`
	Internal string `json:"-"`
}

func (e ErrorStatusNotAllowed) Error() string {
	return e.Message
}

// ErrorStatusNotFound type for not found
type ErrorStatusNotFound struct {
	Message  string `json:"message"`
	Internal string `json:"-"`
}

func (e ErrorStatusNotFound) Error() string {
	return e.Message
}

// ErrorInternalServerError type for internal server error
type ErrorInternalServerError struct {
	Message  string `json:"message"`
	Internal string `json:"-"`
}

func (e ErrorInternalServerError) Error() string {
	return e.Message
}

// ErrorUnprocessableEntity type for validation errors
type ErrorUnprocessableEntity struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Internal string `json:"-"`
}

func (e ErrorUnprocessableEntity) Error() string {
	return e.Message
}
