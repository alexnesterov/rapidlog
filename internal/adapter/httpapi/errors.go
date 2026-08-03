package httpapi

import "errors"

var errInvalidRequest = errors.New("invalid request body")
var errInternal = errors.New("internal server error")
var errParseID = errors.New("parse id error")
