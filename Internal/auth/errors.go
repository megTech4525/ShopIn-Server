package auth

import (
	"errors"
)

var (
	ErrEmailTaken         = errors.New("email already registered")
	ErrPhoneTaken         = errors.New("phone number already registered")
	ErrAccountNotFound    = errors.New("Account Not Found")
	ErrPasswordTooLong    = errors.New("Password too long")
	ErrInvalidCredentials = errors.New("Invalid phone number of password")
	ErrPasswordTooShort   = errors.New("Password too short")
	ErrWeakPassword       = errors.New("Password is weak")
	ErrinvalidPhoneNumber = errors.New("invalid phone number or email")
	ErrInvalidToken       = errors.New("invalid Token")
	ErrTooManyAttempts    = errors.New("too many login attempts, try again later")
)
