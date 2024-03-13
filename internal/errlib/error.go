package errlib

import (
	"errors"
)

// ErrInternal is error to be matched with 500 http code.
var ErrInternal = errors.New("some internal error happened")

// ErrResourceExists is error to be matched with 409 http code.
// Error is raised when you try to create existing resource.
var ErrResourceExists = errors.New("resource already exists")

var ErrNotFound = errors.New("resource not found")
