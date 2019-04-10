package app

import (
	"encoding/json"
	"net/http"

	"github.com/torniker/goapp/app/logger"
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
		ctx.response.WriteHeader(http.StatusBadRequest)
		e := err.(ErrorBadRequest)
		body = ResponseErr{
			Code:    e.Code,
			Message: e.Message,
		}
	case ErrorStatusUnauthorized:
		ctx.response.WriteHeader(http.StatusUnauthorized)
		e := err.(ErrorStatusUnauthorized)
		body = ResponseErr{
			Code:    e.Code,
			Message: e.Message,
		}
	case ErrorStatusNotAllowed:
		ctx.response.WriteHeader(http.StatusMethodNotAllowed)
		e := err.(ErrorStatusNotAllowed)
		body = ResponseErr{
			Code:    e.Code,
			Message: e.Message,
		}
	case ErrorStatusNotFound:
		ctx.response.WriteHeader(http.StatusNotFound)
		e := err.(ErrorStatusNotFound)
		body = ResponseErr{
			Code:    e.Code,
			Message: e.Message,
		}
	default:
		ctx.response.WriteHeader(http.StatusInternalServerError)
		e := err.(ErrorInternalServerError)
		body = ResponseErr{
			Code:    e.Code,
			Message: e.Message,
		}
	}
	logger.ErrorWithCaller(logger.Caller(), err.Error())
	json.NewEncoder(ctx.response).Encode(body)
}

// NotFound response 404
func (ctx *Ctx) NotFound() {
	e := ErrorStatusNotFound{
		Code:    4,
		Message: "not found",
	}
	ctx.Error(e)
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
	return "unauthorized"
}

// ErrorStatusNotAllowed type for not allowed
type ErrorStatusNotAllowed struct {
	Code    int
	Message string
}

func (e ErrorStatusNotAllowed) Error() string {
	return "not allowed"
}

// ErrorStatusNotFound type for not found
type ErrorStatusNotFound struct {
	Code    int
	Message string
}

func (e ErrorStatusNotFound) Error() string {
	return "not found"
}

// ErrorInternalServerError type for internal server error
type ErrorInternalServerError struct {
	Code    int
	Message string
}

func (e ErrorInternalServerError) Error() string {
	return e.Message
}
