package service

import "fmt"

var (
	ErrUserIsNotResposible     = fmt.Errorf("user is not resposible for that organization")
	ErrUserNotExists           = fmt.Errorf("user not exists")
	ErrCannotCreateTender      = fmt.Errorf("cannot create tender")
	ErrTendersNotFound         = fmt.Errorf("tenders not found")
	ErrTenderNotFound          = fmt.Errorf("tender not found")
	ErrBidNotFound             = fmt.Errorf("bid not found")
	ErrCannotGetTenderStatus   = fmt.Errorf("cannot get tender status")
	ErrTenderOrVersionNotFound = fmt.Errorf("tender or version not found")
	ErrBidOrVersionNotFound    = fmt.Errorf("bid or version not found")
)
