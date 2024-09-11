package service

import "errors"

var (
	ErrInvalidAccessToken  = errors.New("INVALID_ACCESS_TOKEN")
	ErrInvalidRefreshToken = errors.New("INVALID_REFRESH_TOKEN")
	ErrUserNotFound        = errors.New("USER_NOT_FOUND")
	ErrSessionNotFound     = errors.New("SESSION_NOT_FOUND")
	ErrInvalidUserID       = errors.New("INVALID_USER_ID")
)
