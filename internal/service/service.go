package service

import (
	"github.com/google/uuid"
	forum "stepik.leoscode.http/internal/gen/api"
	"stepik.leoscode.http/internal/models"
	"stepik.leoscode.http/internal/repository"
)

type Repo interface { // изменить интерфейс под дальнейшую реализацию
	Authorize(username, pwd string) uuid.UUID
	CheckUsernameExist(username string) (uuid.UUID, error)
	CheckCorrectPwd(username, password string) error

	CreateThread(thread forum.ThreadCreate, XUXI repository.XUXI) (forum.Thread, error)
	GetThread(id forum.ThreadIdPath) (forum.Thread, error)
	GetThreads(filter repository.ThreadListFilter) (forum.ThreadListResponse, error)
	ReplaceThreadById(id forum.ThreadIdPath, create forum.ThreadCreate) (forum.Thread, error)
	ChangeThreadById(id forum.ThreadIdPath, patch models.ThreadPatchInput) (forum.Thread, error)
	DeleteThreadByUd(id forum.ThreadIdPath) error

	CreatePost(post forum.PostCreate, id forum.ThreadIdPath, XUXI repository.XUXI) (forum.Post, error)
	CheckUserExist(user uuid.UUID) error
	CheckThreadAlreadyExist(XUXI repository.XUXI) (forum.Thread, error)
	CheckPostAlreadyExist(XUXI repository.XUXI) (forum.Post, error)
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
