package errors

import "errors"

var (
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrFailedToGenerateJWT  = errors.New("failed to generate jwt")
	ErrFailedToHashPassword = errors.New("failed to hash password")
	ErrFailedToAddHabit     = errors.New("failed to add habit")
)
