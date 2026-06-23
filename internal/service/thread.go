package service

import (
	"errors"
	"reflect"

	"github.com/google/uuid"
	forum "stepik.leoscode.http/internal/gen/api"
	"stepik.leoscode.http/internal/repository"
)

var (
	ErrAlreadyThreadExist = errors.New("thread already exist")
	ErrConflict           = errors.New("conflict")
	ErrNoSuchUserExist    = errors.New("no such user exist")

	ErrThreadNotFound = errors.New("thread not found")
)

func (s *Service) CreateThread(threadCreate forum.ThreadCreate, XUXI repository.XUXI) (forum.Thread, error) {
	if err := s.repo.CheckUserExist(XUXI.XU); errors.Is(err, repository.ErrNoSuchUserExist) {
		return forum.Thread{}, ErrNoSuchUserExist
	}
	if t, err := s.repo.CheckThreadAlreadyExist(XUXI); errors.Is(err, repository.ErrUserIdAlreadyExist) {
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

func (s *Service) GetThreadById(id int64, userId *uuid.UUID) (forum.Thread, error) {
	if userId != nil {
		err := s.repo.CheckUserExist(*userId)
		if errors.Is(err, repository.ErrNoSuchUserExist) {
			return forum.Thread{}, ErrNoSuchUserExist
		}
	}

	t, err := s.repo.GetThread(id)
	if errors.Is(err, repository.ErrNoThreadFound) {
		return forum.Thread{}, ErrThreadNotFound
	}

	return t, nil
}

func (s *Service) GetListThreads(threadFilter repository.ThreadListFilter) (forum.ThreadListResponse, error) {

	return s.repo.GetThreads(threadFilter)
}
