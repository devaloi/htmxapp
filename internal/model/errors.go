package model

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrDuplicateEmail = errors.New("email already exists")
)
