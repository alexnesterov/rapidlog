package entity

import "errors"

type ValidationError struct {
	Err error
}

func (e *ValidationError) Error() string { return e.Err.Error() }
func (e *ValidationError) Unwrap() error { return e.Err }

var (
	ErrNotFound = errors.New("not found")
)
