package constants

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrEmailExists     = errors.New("email already exists")
	ErrUserDeactivated = errors.New("user account is deactivated")
	ErrRecordNotFound  = errors.New("record not found")
)
