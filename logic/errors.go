package logic

import "errors"

var (
	ErrNotFound     = errors.New("record not found")
	ErrDuplicate    = errors.New("duplicate record")
	ErrUnauthorised = errors.New("unauthorised")
	ErrPostTooOld   = errors.New("post is too old to vote on")

	// User errors
	ErrUsernameTaken      = errors.New("username already taken")
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrAccountSuspended   = errors.New("account suspended")

	// Community errors
	ErrCommunityNotFound = errors.New("community not found")
)
