package repository

import "errors"

var (
	ErrUserNotFound    = errors.New("USER_NOT_FOUND")
	ErrSessionNotFound = errors.New("SESSION_NOT_FOUND")
)
