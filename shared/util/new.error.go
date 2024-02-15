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
	SessionNotFound         = errors.New("Session not found")
	SessionExpired          = errors.New("Session expired")
	SessionIsBlocked        = errors.New("Session is blocked")
	SessionNotMatchUser     = errors.New("Session not match user")
	InvalidRefreshToken     = errors.New("Invalid refresh token")
	RefCodeNotFound         = errors.New("RefCode not Found")
	UserRegisterNotFound    = errors.New("User Register not Found")
)
