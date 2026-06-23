package inMemoryRepo

import (
	"time"

	"github.com/google/uuid"
	forum "stepik.leoscode.http/internal/gen/api"
	"stepik.leoscode.http/internal/repository"
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

func (r *RepoInMemory) GetThread(id int64) (forum.Thread, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if t, ok := r.idToThread[id]; ok {
		return t, nil
	}
	return forum.Thread{}, repository.ErrNoThreadFound
}

func (r *RepoInMemory) GenerateThreadId() int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.threadId++
	return r.threadId
}
func (r *RepoInMemory) CheckUserExist(user uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.UUIDToPwd[user]
	if !ok {
		return repository.ErrNoSuchUserExist
	}
	return nil
}
func (r *RepoInMemory) CheckThreadAlreadyExist(XUXI repository.XUXI) (forum.Thread, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	thread, ok := r.XUXIToThread[XUXI]
	if ok {
		return thread, repository.ErrUserIdAlreadyExist
	}
	return forum.Thread{}, nil
}
