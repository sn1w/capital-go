package error

import "errors"

var (
	ErrUnknown               = errors.New("unknown")
	ErrBadRequest            = errors.New("bad request")
	ErrUnAuthorized          = errors.New("unauthorized")
	ErrResourceNotFound      = errors.New("resource not found")
	ErrUnknownResponseFormat = errors.New("unknown resposne format")
)
