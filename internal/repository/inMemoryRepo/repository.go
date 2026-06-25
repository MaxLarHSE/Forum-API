package inMemoryRepo

import (
	"sync"

	"github.com/google/uuid"
	forum "stepik.leoscode.http/internal/gen/api"
	"stepik.leoscode.http/internal/repository"
	"stepik.leoscode.http/internal/service"
)

var _ service.Repo = &RepoInMemory{}

type RepoInMemory struct {
	userToUUID map[string]uuid.UUID
	UUIDToPwd  map[uuid.UUID]string

	idToThread   map[int64]forum.Thread // тут лучше модели форума или инт64
	XUXIToThread map[repository.XUXI]forum.Thread

	idToPost   map[int64]forum.Post
	XUXIToPost map[repository.XUXI]forum.Post

	idThreadToIdPost map[int64][]int64
	mu               sync.Mutex
	threadId         int64
	postId           int64
}

func (r *RepoInMemory) Clear() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	clear(r.userToUUID)
	clear(r.UUIDToPwd)
	clear(r.idToThread)
	clear(r.XUXIToThread)
	clear(r.idToPost)
	clear(r.XUXIToPost)
	clear(r.idThreadToIdPost)

	return nil
}

func NewRepoInMemory() *RepoInMemory {
	return &RepoInMemory{
		userToUUID:       make(map[string]uuid.UUID),
		UUIDToPwd:        make(map[uuid.UUID]string),
		idToThread:       make(map[int64]forum.Thread),
		XUXIToThread:     make(map[repository.XUXI]forum.Thread),
		idToPost:         make(map[int64]forum.Post),
		XUXIToPost:       make(map[repository.XUXI]forum.Post),
		idThreadToIdPost: make(map[int64][]int64),
	}
}
