package httpapi

import "errors"

var ErrInvalidRequestBody = errors.New("invalid request body")
var ErrInvalidID = errors.New("invalid id")
var ErrInternal = errors.New("internal server error")
