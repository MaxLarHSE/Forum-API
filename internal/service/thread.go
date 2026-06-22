package service

import (
	"errors"
	"reflect"

	forum "stepik.leoscode.http/internal/gen/api"
	"stepik.leoscode.http/internal/repository"
)

var (
	ErrAlreadyThreadExist = errors.New("thread already exist")
	ErrConflict           = errors.New("conflict")
	ErrNoSuchUserExist    = errors.New("no such user exist")
)

func (s *Service) CreateThread(threadCreate forum.ThreadCreate, XUXI repository.XUXI) (forum.Thread, error) {
	if !s.repo.CheckUserExist(XUXI.XU) {
		return forum.Thread{}, ErrNoSuchUserExist
	}
	if t, ok := s.repo.CheckThreadAlreadyExist(XUXI); ok {
		var tc = forum.ThreadCreate{
			Content: t.Content,
			Tags:    t.Tags,
			Title:   t.Title,
		}
		if !reflect.DeepEqual(tc, threadCreate) {
			return forum.Thread{}, ErrConflict
		}
		return t, ErrAlreadyThreadExist
	}

	return s.repo.CreateThread(threadCreate, XUXI), nil
}
