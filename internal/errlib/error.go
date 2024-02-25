package errlib

import (
	"errors"
)

// ErrHttpInternal is error to be matched with 500 http code
var ErrHttpInternal = errors.New("some internal error happened")
