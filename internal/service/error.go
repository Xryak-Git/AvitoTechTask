package service

import "fmt"

var (
	ErrUserIsNotResposible = fmt.Errorf("user is not resposible for that organization")
	ErrUserNotExists       = fmt.Errorf("user not exists")
	ErrCannotCreateTender  = fmt.Errorf("cannot create tender")
)
