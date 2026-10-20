package auth

import (
	"errors"
)

var (
	ErrEmailTaken      = errors.New("email already registered")
	ErrPhoneTaken      = errors.New("phone number already registered")
	ErrAccountNotFound = errors.New("Account Not Found")
	ErrPasswordTooLong = errors.New("Password too long")
)
