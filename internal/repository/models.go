package repository

import (
	"errors"

	"github.com/google/uuid"
)

type XUXI struct { // где лучше разместить
	XU uuid.UUID
	XI string
}

var ( // переписать ошибки
	ErrAlreadyThreadExist = errors.New("thread already exist")
	ErrConflict           = errors.New("conflict")
	ErrNoSuchUserExist    = errors.New("no such user exist")
	ErrNoThreadFound      = errors.New("thread not exist")

	ErrUserNotExist       = errors.New("user not exist")
	ErrUserIdAlreadyExist = errors.New("user id already exist")
	ErrPwdNotCorrect      = errors.New("password not correct")
	ErrUserAlreadyExist   = errors.New("user already exist")
)
