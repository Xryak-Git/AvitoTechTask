package repoerrs

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrNotExists     = errors.New("not exists")
	ErrAlreadyExists = errors.New("already exists")
	ErrURLNotFound   = fmt.Errorf("url not found")
)
