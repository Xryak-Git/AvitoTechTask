package storage

import "fmt"

var (
	ErrURLNotFound      = fmt.Errorf("url not found")
	ErrURLAlreadyExists = fmt.Errorf("url already exists")
)
