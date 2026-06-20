package service

import (
	"errors"

	"github.com/google/uuid"
)

type Repo interface {
	Authorize(username, pwd string) uuid.UUID
	CheckUsernameExist(username string) (uuid.UUID, bool)
	CheckCorrectPwd(username, password string) bool
	Clear()
}
type Service struct {
	repo Repo
}

func NewService(r Repo) *Service {
	return &Service{repo: r}
}

var (
	InvalidPasswordError = errors.New("Invalid password")
)

func (s *Service) Authorize(username string, pwd string) (uuid.UUID, error) {
	userUUID, exist := s.repo.CheckUsernameExist(username)
	if exist {
		if !s.repo.CheckCorrectPwd(username, pwd) {
			return uuid.UUID{}, InvalidPasswordError
		}
		return userUUID, nil
	}

	return s.repo.Authorize(username, pwd), nil
}

func (s *Service) Truncate() {
	s.repo.Clear()
}
