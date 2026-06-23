package inMemoryRepo

import (
	"github.com/google/uuid"
	"stepik.leoscode.http/internal/repository"
)

func (r *RepoInMemory) Authorize(username, pwd string) uuid.UUID {
	newId := uuid.New()
	r.userToUUID[username] = newId
	r.UUIDToPwd[newId] = pwd
	return newId
}

func (r *RepoInMemory) CheckUsernameExist(username string) (uuid.UUID, error) {
	id, exist := r.userToUUID[username]
	if exist {
		return id, repository.ErrUserAlreadyExist
	}
	return id, nil
}
func (r *RepoInMemory) CheckCorrectPwd(username, password string) error {
	if r.UUIDToPwd[r.userToUUID[username]] == password {
		return nil
	}
	return repository.ErrPwdNotCorrect
}
