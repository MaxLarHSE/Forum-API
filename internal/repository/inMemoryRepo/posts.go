package inMemoryRepo

import (
	"log"
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
	r.idThreadToIdPost[threadId] = append(r.idThreadToIdPost[threadId], id)
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
func (r *RepoInMemory) GetPosts(id forum.ThreadIdPath, filter repository.PostListFilter) (forum.PostListResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	postsId := r.idThreadToIdPost[id]
	res := make([]forum.Post, 0)
	for i := filter.Offset; i < (int32)(len(postsId)) && i < filter.Offset+filter.Limit; i++ {
		res = append(res, r.idToPost[postsId[i]])
	}
	log.Println(r.idToPost, r.idThreadToIdPost, filter)
	return forum.PostListResponse{
		Items: res,
		Meta: forum.PaginationMeta{
			Limit:  filter.Limit,
			Offset: filter.Offset,
			Total:  (int64)(len(postsId)),
		},
	}, nil
}
