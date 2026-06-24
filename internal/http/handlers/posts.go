package handlers

import (
	"errors"
	"net/http"

	forum "stepik.leoscode.http/internal/gen/api"
	"stepik.leoscode.http/internal/repository"
	"stepik.leoscode.http/internal/service"
)

func (s *Server) ListPosts(w http.ResponseWriter, r *http.Request, threadId forum.ThreadIdPath, params forum.ListPostsParams) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) CreatePost(w http.ResponseWriter, r *http.Request, threadId forum.ThreadIdPath, params forum.CreatePostParams) {
	//Post
	if !validateThreadId(threadId) {
		writeError(w, http.StatusBadRequest, forum.ErrorBadRequest{forum.ValidationError, &map[string]any{"thread id": threadId}, "thread is minimum 1"})
		return
	}
	if !validateXI(params.XIdempotencyKey) {
		writeError(w, http.StatusBadRequest, forum.ErrorUnauthorized{
			forum.ValidationError,
			&map[string]any{"XIdempotencyKey": params.XIdempotencyKey},
			"err length XIdempotencyKey"},
		)
		return
	}
	var postCreate forum.PostCreate
	if err := readJson(r, &postCreate); err != nil {
		writeError(w, http.StatusBadRequest, forum.ErrorBadRequest{
			forum.BadRequest,
			nil,
			"not json content type"},
		)
		return
	}
	post, err := s.service.CreatePost(postCreate, threadId, repository.XUXI{
		XU: params.XUserId,
		XI: params.XIdempotencyKey,
	})
	switch {
	case err == nil:
		writeJson(w, http.StatusCreated, post)
	case errors.Is(err, service.ErrConflict):
		writeError(w, http.StatusConflict, forum.ErrorConflict{
			forum.BadRequest,
			nil,
			"distinct body but equal user and idempotency-key params in memory",
		})
	case errors.Is(err, service.ErrThreadNotFound):
		writeError(w, http.StatusNotFound, forum.ErrorNotFound{
			forum.NotFound,
			nil,
			"no such thread exist",
		})
	case errors.Is(err, service.ErrAlreadyPostExist):
		writeJson(w, http.StatusOK, post)
	case errors.Is(err, service.ErrUserNotExist):
		writeError(w, http.StatusUnauthorized, forum.ErrorUnauthorized{
			forum.InvalidCredentials,
			nil,
			"no such user exist",
		})
	case errors.Is(err, service.ErrTryChangeLockedThread):
		writeError(w, http.StatusForbidden, forum.ErrorForbidden{
			forum.Forbidden,
			nil,
			"thy change locked thread",
		})
	default:
		writeError(w, http.StatusInternalServerError, forum.ErrorInternal{
			forum.InternalError,
			&map[string]interface{}{"error": err.Error()},
			"no relevant error",
		})
	}
}
