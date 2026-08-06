package httpapi

import "errors"

var errInvalidRequestBody = errors.New("invalid request body")
var errInvalidID = errors.New("invalid id")
var errInternal = errors.New("internal server error")
