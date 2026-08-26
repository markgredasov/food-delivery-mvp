package errs

import "errors"

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrNotFound        = errors.New("not found")
	ErrAlreadyExists   = errors.New("already exists")
	ErrInvalidRequest  = errors.New("invalid request")
	ErrNotImplemented  = errors.New("not implemented yet")
)
