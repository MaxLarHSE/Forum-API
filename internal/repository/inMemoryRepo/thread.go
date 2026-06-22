package inMemoryRepo

import (
	"errors"
	"time"

	"github.com/google/uuid"
	forum "stepik.leoscode.http/internal/gen/api"
	"stepik.leoscode.http/internal/repository"
)

var (
	ErrAlreadyThreadExist = errors.New("thread already exist")
	ErrConflict           = errors.New("conflict")
	ErrNoSuchUserExist    = errors.New("no such user exist")
)

func (r *RepoInMemory) CreateThread(threadCreate forum.ThreadCreate, XUXI repository.XUXI) forum.Thread {
	id := r.GenerateThreadId()
	thread := forum.Thread{
		AuthorId:  XUXI.XU,
		Content:   threadCreate.Content,
		CreatedAt: time.Now(),
		Id:        id,
		IsLocked:  false,
		Tags:      threadCreate.Tags,
		Title:     threadCreate.Title,
		UpdatedAt: nil,
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.XUXIToThread[XUXI] = thread
	r.idToThread[id] = thread
	return thread
}
func (r *RepoInMemory) GenerateThreadId() int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.threadId++
	return r.threadId
}
func (r *RepoInMemory) CheckUserExist(user uuid.UUID) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.UUIDToPwd[user]
	return ok
}
func (r *RepoInMemory) CheckThreadAlreadyExist(XUXI repository.XUXI) (forum.Thread, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	thread, ok := r.XUXIToThread[XUXI]
	return thread, ok
}
