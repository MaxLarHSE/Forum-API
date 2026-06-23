package service

import (
	"errors"

	"github.com/google/uuid"
	"stepik.leoscode.http/internal/repository"
)

var (
	ErrInvalidPassword = errors.New("Invalid password")
)

func (s *Service) Authorize(username string, pwd string) (uuid.UUID, error) {
	userUUID, err := s.repo.CheckUsernameExist(username)

	if errors.Is(err, repository.ErrUserAlreadyExist) {
		if err = s.repo.CheckCorrectPwd(username, pwd); errors.Is(err, repository.ErrPwdNotCorrect) {
			return uuid.UUID{}, ErrInvalidPassword
		}
		return userUUID, nil
	}

	return s.repo.Authorize(username, pwd), nil
}
