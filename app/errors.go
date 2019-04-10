package app

import (
	"net/http"
)

// ResponseErr describes error response
type ResponseErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error checks error type and responses accordingly
func (ctx *Ctx) Error(err error) {
	var body ResponseErr
	switch err.(type) {
	case ErrorBadRequest:
		ctx.response.SetStatus(http.StatusBadRequest)
		e := err.(ErrorBadRequest)
		body = ResponseErr{
			Code:    e.Code,
			Message: e.Message,
		}
	case ErrorStatusUnauthorized:
		ctx.response.SetStatus(http.StatusUnauthorized)
		e := err.(ErrorStatusUnauthorized)
		body = ResponseErr{
			Code:    e.Code,
			Message: e.Message,
		}
	case ErrorStatusNotAllowed:
		ctx.response.SetStatus(http.StatusMethodNotAllowed)
		e := err.(ErrorStatusNotAllowed)
		body = ResponseErr{
			Code:    e.Code,
			Message: e.Message,
		}
	case ErrorStatusNotFound:
		ctx.response.SetStatus(http.StatusNotFound)
		e := err.(ErrorStatusNotFound)
		body = ResponseErr{
			Code:    e.Code,
			Message: e.Message,
		}
	case ErrorInternalServerError:
		ctx.response.SetStatus(http.StatusInternalServerError)
		e := err.(ErrorInternalServerError)
		body = ResponseErr{
			Code:    e.Code,
			Message: e.Message,
		}
	default:
		ctx.response.SetStatus(http.StatusInternalServerError)
		body = ResponseErr{
			Code:    99,
			Message: err.Error(),
		}
	}
	ctx.response.Write(body)
}

// NotFound response 404
func (ctx *Ctx) NotFound() error {
	e := ErrorStatusNotFound{
		Code:    4,
		Message: "not found",
	}
	return e
}

// Unauthorized response 401
func (ctx *Ctx) Unauthorized() error {
	e := ErrorStatusUnauthorized{
		Code:    1,
		Message: "unauthorizes",
	}
	return e
}

// InternalError response 404
func (ctx *Ctx) InternalError() error {
	e := ErrorInternalServerError{
		Code:    5,
		Message: "internal server error",
	}
	return e
}

// ErrorBadRequest type for bad request
type ErrorBadRequest struct {
	Code    int
	Message string
}

func (e ErrorBadRequest) Error() string {
	return e.Message
}

// ErrorStatusUnauthorized type for Unauthorized
type ErrorStatusUnauthorized struct {
	Code    int
	Message string
}

func (e ErrorStatusUnauthorized) Error() string {
	return e.Message
}

// ErrorStatusNotAllowed type for not allowed
type ErrorStatusNotAllowed struct {
	Code    int
	Message string
}

func (e ErrorStatusNotAllowed) Error() string {
	return e.Message
}

// ErrorStatusNotFound type for not found
type ErrorStatusNotFound struct {
	Code    int
	Message string
}

func (e ErrorStatusNotFound) Error() string {
	return e.Message
}

// ErrorInternalServerError type for internal server error
type ErrorInternalServerError struct {
	Code    int
	Message string
}

func (e ErrorInternalServerError) Error() string {
	return e.Message
}
