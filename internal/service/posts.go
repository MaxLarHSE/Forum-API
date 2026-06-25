package service

import (
	"errors"

	forum "stepik.leoscode.http/internal/gen/api"
	"stepik.leoscode.http/internal/repository"
)

var (
	ErrAlreadyPostExist         = errors.New("post already exist")
	ErrTryAddPostToLockedThread = errors.New("thy add post to locked thread")
)

func (s *Service) CreatePost(postCreate forum.PostCreate, id forum.ThreadIdPath, XUXI repository.XUXI) (forum.Post, error) {
	if err := s.repo.CheckUserExist(XUXI.XU); errors.Is(err, repository.ErrUserNotExist) {
		return forum.Post{}, ErrUserNotExist
	}
	thread, err := s.repo.GetThread(id)
	if errors.Is(err, repository.ErrNoThreadFound) {
		return forum.Post{}, ErrThreadNotFound
	}
	if thread.IsLocked {
		return forum.Post{}, ErrTryAddPostToLockedThread
	}
	if p, err := s.repo.CheckPostAlreadyExist(XUXI); errors.Is(err, repository.ErrUserIdAlreadyExist) {
		if p.Content != postCreate.Content {
			return forum.Post{}, ErrConflict
		}
		return p, ErrAlreadyPostExist
	}

	return s.repo.CreatePost(postCreate, id, XUXI)
}

func (s *Service) GetListPosts(id forum.ThreadIdPath, filter repository.PostListFilter) (forum.PostListResponse, error) {
	_, err := s.repo.GetThread(id)
	if errors.Is(err, repository.ErrNoThreadFound) {
		return forum.PostListResponse{}, ErrThreadNotFound
	}
	list, err := s.repo.GetPosts(id, filter)
	return list, err
}
