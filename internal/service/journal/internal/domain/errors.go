package domain

import "errors"

type ValidationError struct {
	Err error
}

func (e *ValidationError) Error() string { return e.Err.Error() }
func (e *ValidationError) Unwrap() error { return e.Err }

var (
	ErrOwnerMismatch = errors.New("owned by another user")
	ErrNotFound      = errors.New("not found")
)
