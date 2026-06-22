package service

import (
	"github.com/google/uuid"
	forum "stepik.leoscode.http/internal/gen/api"
	"stepik.leoscode.http/internal/repository"
)

type Repo interface {
	Authorize(username, pwd string) uuid.UUID
	CheckUsernameExist(username string) (uuid.UUID, bool)
	CheckCorrectPwd(username, password string) bool

	CreateThread(thread forum.ThreadCreate, XUXI repository.XUXI) forum.Thread
	CheckUserExist(user uuid.UUID) bool
	CheckThreadAlreadyExist(XUXI repository.XUXI) (forum.Thread, bool)
	Clear() error
}
type Service struct {
	repo Repo
}

func NewService(r Repo) *Service {
	return &Service{repo: r}
}

func (s *Service) Truncate() error {
	return s.repo.Clear()
}
