package service

import (
	"errors"
	"reflect"

	"github.com/google/uuid"
	forum "stepik.leoscode.http/internal/gen/api"
	"stepik.leoscode.http/internal/models"
	"stepik.leoscode.http/internal/repository"
)

var (
	ErrAlreadyThreadExist = errors.New("thread already exist")
	ErrConflict           = errors.New("conflict")
	ErrUserNotExist       = errors.New("no such user exist")

	ErrThreadNotFound        = errors.New("thread not found")
	ErrUserDontHaveRights    = errors.New("user dont have rights")
	ErrTryChangeLockedThread = errors.New("try change locked thread")
)

func (s *Service) CreateThread(threadCreate forum.ThreadCreate, XUXI repository.XUXI) (forum.Thread, error) {
	if err := s.repo.CheckUserExist(XUXI.XU); errors.Is(err, repository.ErrUserNotExist) {
		return forum.Thread{}, ErrUserNotExist
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
		if errors.Is(err, repository.ErrUserNotExist) {
			return forum.Thread{}, ErrUserNotExist
		}
	}

	t, err := s.repo.GetThread(id)
	if errors.Is(err, repository.ErrNoThreadFound) {
		return forum.Thread{}, ErrThreadNotFound
	}

	return t, nil
}

func (s *Service) ReplaceThreadById(id forum.ThreadIdPath, create forum.ThreadCreate, params forum.ReplaceThreadParams) (forum.Thread, error) {
	thread, err := s.repo.GetThread(id)
	if errors.Is(err, repository.ErrNoThreadFound) {
		return forum.Thread{}, ErrThreadNotFound
	}
	err = s.repo.CheckUserExist(params.XUserId)
	if errors.Is(err, repository.ErrUserNotExist) {
		return forum.Thread{}, ErrUserNotExist
	}
	if thread.AuthorId != params.XUserId { // в бд или тут?
		return forum.Thread{}, ErrUserDontHaveRights
	}
	newThread, err := s.repo.ReplaceThreadById(id, create)
	return newThread, err
}

func (s *Service) ChangeThreadById(id forum.ThreadIdPath, patch models.ThreadPatchInput, params forum.PatchThreadParams) (forum.Thread, error) {
	thread, err := s.repo.GetThread(id)
	if errors.Is(err, repository.ErrNoThreadFound) {
		return forum.Thread{}, ErrThreadNotFound
	}
	if thread.IsLocked && !(patch.IsLocked != nil && *patch.IsLocked == false) {

		return forum.Thread{}, ErrTryChangeLockedThread
	}
	err = s.repo.CheckUserExist(params.XUserId)
	if errors.Is(err, repository.ErrUserNotExist) {
		return forum.Thread{}, ErrUserNotExist
	}
	if thread.AuthorId != params.XUserId { // в бд или тут?
		return forum.Thread{}, ErrUserDontHaveRights
	}

	changedThread, err := s.repo.ChangeThreadById(id, patch)
	return changedThread, err
}

func (s *Service) GetListThreads(threadFilter repository.ThreadListFilter) (forum.ThreadListResponse, error) {

	return s.repo.GetThreads(threadFilter)
}
func (s *Service) DeleteThread(id forum.ThreadIdPath, params forum.DeleteThreadParams) error {
	thread, err := s.repo.GetThread(id)
	if errors.Is(err, repository.ErrNoThreadFound) {
		return ErrThreadNotFound
	}
	err = s.repo.CheckUserExist(params.XUserId)
	if errors.Is(err, repository.ErrUserNotExist) {
		return ErrUserNotExist
	}
	if thread.AuthorId != params.XUserId { // в бд или тут?
		return ErrUserDontHaveRights
	}
	return s.repo.DeleteThreadByUd(id)
}
