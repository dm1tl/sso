package storage

import "errors"

var (
	ErrNotExists     = errors.New("user not exists")
	ErrAlreadyExists = errors.New("user already exists")
	ErrAppNotExists  = errors.New("app is not exists")
)
