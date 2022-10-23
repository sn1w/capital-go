package error

import "errors"

var (
	ErrUnknown          = errors.New("unknown")
	ErrUnAuthorized     = errors.New("unauthorized")
	ErrResourceNotFound = errors.New("resource not found")
)
