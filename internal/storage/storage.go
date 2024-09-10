package storage

import "fmt"

var (
	ErrURLNotFound   = fmt.Errorf("url not found")
	ErrAlreadyExists = fmt.Errorf("url already exists")
)
