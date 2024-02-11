package util

import "errors"

var (
	EmailAlreadyExist       = errors.New("Email already exists")
	EmailAlreadyUsed        = errors.New("Email already used")
	NimAlreadyUsed          = errors.New("Nim already used")
	EmailNotFound           = errors.New("Email not found")
	InvalidPassword         = errors.New("Invalid password")
	UsernameOrEmailNotFound = errors.New("Username or email not found")
	UserNotFound            = errors.New("User not found")
)
