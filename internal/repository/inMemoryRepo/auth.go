package inMemoryRepo

import "github.com/google/uuid"

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
func (r *RepoInMemory) CheckCorrectPwd(username, password string) bool {
	return r.UUIDToPwd[r.userToUUID[username]] == password
}
