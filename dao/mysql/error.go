package mysql

import "errors"

// Sentinel errors — logic layer checks these to decide HTTP status codes.
var (
	ErrNotFound     = errors.New("record not found")
	ErrDuplicate    = errors.New("duplicate record")
	ErrUnauthorised = errors.New("unauthorised")
)
