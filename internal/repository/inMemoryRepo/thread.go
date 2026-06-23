package inMemoryRepo

import (
	"slices"
	"sort"
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

func (r *RepoInMemory) GetThreads(filter repository.ThreadListFilter) (forum.ThreadListResponse, error) {
	threads := r.ThreadsSortedByID()
	res := []forum.Thread{}
	isProper := func(t forum.Thread) bool {
		if filter.Tag != nil {
			if t.Tags == nil || !slices.Contains(*t.Tags, *filter.Tag) {
				return false
			}
		}
		if filter.AuthorID != nil {
			if t.AuthorId != *filter.AuthorID {
				return false
			}
		}
		return true
	}
	for i := range threads {
		if isProper(threads[i]) {
			res = append(res, threads[i])
		}
	}
	var total = (int64)(len(res))
	end := min((int)(filter.Offset+filter.Limit), len(res))
	begin := min((int)(filter.Offset), len(res))
	result := make([]forum.Thread, end-begin)
	copy(result, res[begin:end])
	return forum.ThreadListResponse{Items: result, Meta: forum.PaginationMeta{
		Limit:  filter.Limit,
		Offset: filter.Offset,
		Total:  total,
	}}, nil
}

func (r *RepoInMemory) ThreadsSortedByID() []forum.Thread { // gpt func
	r.mu.Lock()
	defer r.mu.Unlock()

	ids := make([]int64, 0, len(r.idToThread))
	for id := range r.idToThread {
		ids = append(ids, id)
	}

	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})

	threads := make([]forum.Thread, 0, len(ids))
	for _, id := range ids {
		threads = append(threads, r.idToThread[id])
	}

	return threads
}
