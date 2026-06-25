package repository

import (
	"errors"

	"github.com/google/uuid"
)

// где лучше разместить
type XUXI struct {
	XU uuid.UUID
	XI string
}

var ( // переписать ошибки
	ErrAlreadyThreadExist = errors.New("thread already exist")
	ErrConflict           = errors.New("conflict")
	ErrUserNotExist       = errors.New("user no exist")
	ErrNoThreadFound      = errors.New("thread not exist")

	ErrUserIdAlreadyExist = errors.New("user id already exist")
	ErrPwdNotCorrect      = errors.New("password not correct")
	ErrUserAlreadyExist   = errors.New("user already exist")

	ErrUserDontHaveRights    = errors.New("user dont have rights")
	ErrTryChangeLockedThread = errors.New("try change locked thread")
)

type ThreadListFilter struct {
	Limit    int32
	Offset   int32
	Tag      *string
	AuthorID *uuid.UUID
}
type PostListFilter struct {
	Limit  int32
	Offset int32
}
