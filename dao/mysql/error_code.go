package mysql

import "errors"

var (
	ErrorUserExist       = errors.New("user already exists")
	ErrorUserNotExist    = errors.New("user not exists")
	ErrorInvalidPassword = errors.New("username or password error")
	ErrorInvalidID       = errors.New("invalid id")
)
