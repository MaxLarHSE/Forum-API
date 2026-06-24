package inMemoryRepo

import (
	"time"

	forum "stepik.leoscode.http/internal/gen/api"
	"stepik.leoscode.http/internal/repository"
)

func (r *RepoInMemory) CreatePost(postCreate forum.PostCreate, threadId forum.ThreadIdPath, XUXI repository.XUXI) (forum.Post, error) {
	id := r.GeneratePostId()
	post := forum.Post{
		AuthorId:  XUXI.XU,
		Content:   postCreate.Content,
		CreatedAt: time.Now(),
		Id:        id,
		ThreadId:  threadId,
		UpdatedAt: nil,
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.XUXIToPost[XUXI] = post
	r.idToPost[id] = post
	return post, nil
}
func (r *RepoInMemory) GeneratePostId() int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.postId++
	return r.postId
}
func (r *RepoInMemory) CheckPostAlreadyExist(XUXI repository.XUXI) (forum.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	post, ok := r.XUXIToPost[XUXI]
	if ok {
		return post, repository.ErrUserIdAlreadyExist
	}
	return forum.Post{}, nil
}
