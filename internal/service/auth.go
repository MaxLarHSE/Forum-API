package service

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidPassword = errors.New("Invalid password")
)

func (s *Service) Authorize(username string, pwd string) (uuid.UUID, error) {
	userUUID, exist := s.repo.CheckUsernameExist(username)
	if exist {
		if !s.repo.CheckCorrectPwd(username, pwd) {
			return uuid.UUID{}, ErrInvalidPassword
		}
		return userUUID, nil
	}

	return s.repo.Authorize(username, pwd), nil
}
