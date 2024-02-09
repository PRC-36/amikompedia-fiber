package util

import "errors"

var (
	EmailAlreadyExist = errors.New("Email already exists")
	EmailNotFound     = errors.New("Email not found")
	InvalidPassword   = errors.New("Invalid password")
)
