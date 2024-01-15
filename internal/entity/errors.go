package entity

import "errors"

var (
	ErrPersonNotExists = errors.New("person doesn't exist")
	ErrInternalService = errors.New("unexpected server error")
)
