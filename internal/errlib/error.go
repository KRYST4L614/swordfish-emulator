package errlib

import (
	"errors"
	"net/http"
)

// ErrInternal is error to be matched with 500 http code.
var ErrInternal = errors.New("some internal error happened")

// ErrBadRequest is error to be matched with 400 http code.
var ErrBadRequest = errors.New("bad request")

// ErrResourceAlreadyExists is error to be matched with 409 http code.
// Error is raised when you try to create existing resource.
var ErrResourceAlreadyExists = errors.New("resource already exists")

// ErrNotFound is error to be matched with 404 http code.
var ErrNotFound = errors.New("resource not found")

type JSONError struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	} `json:"error"`
}

func GetJSONError(err error) *JSONError {
	var jsonErr = JSONError{}
	jsonErr.Error.Message = "Error: " + err.Error()
	switch {
	case errors.Is(err, ErrInternal):
		jsonErr.Error.Code = http.StatusInternalServerError
	case errors.Is(err, ErrResourceAlreadyExists):
		jsonErr.Error.Code = http.StatusConflict
	case errors.Is(err, ErrNotFound):
		jsonErr.Error.Code = http.StatusNotFound
	case errors.Is(err, ErrBadRequest):
		jsonErr.Error.Code = http.StatusBadRequest
	default:
		jsonErr.Error.Code = http.StatusInternalServerError
	}

	return &jsonErr
}
