package repository

import (
	"github.com/google/uuid"
	"stepik.leoscode.http/internal/service"
)

var _ service.Repo = &RepoInMemory{}

type RepoInMemory struct {
	userToUUID map[string]uuid.UUID
	UUIDToPwd  map[uuid.UUID]string
}

func (r *RepoInMemory) CheckCorrectPwd(username, password string) bool {
	return r.UUIDToPwd[r.userToUUID[username]] == password
}

func NewRepoInMemory() *RepoInMemory {
	return &RepoInMemory{
		userToUUID: make(map[string]uuid.UUID),
		UUIDToPwd:  make(map[uuid.UUID]string),
	}
}
func (r *RepoInMemory) Clear() {
	clear(r.userToUUID)
	clear(r.UUIDToPwd)
}

func (r *RepoInMemory) Authorize(username, pwd string) uuid.UUID {
	newId := uuid.New()
	r.userToUUID[username] = newId
	r.UUIDToPwd[newId] = pwd
	return newId
}

func (r *RepoInMemory) CheckUsernameExist(username string) (uuid.UUID, bool) {
	id, exist := r.userToUUID[username]
	return id, exist
}
