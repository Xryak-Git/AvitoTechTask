package repoerrs

import (
	"errors"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrNotExists     = errors.New("not exists")
	ErrAlreadyExists = errors.New("already exists")
)
