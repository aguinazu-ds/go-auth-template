package db

import "errors"

var (
	ErrDuplicateEmail = errors.New("el email ya está en uso")
)
