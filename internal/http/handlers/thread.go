package handlers

import (
	"errors"
	"net/http"

	forum "stepik.leoscode.http/internal/gen/api"
	"stepik.leoscode.http/internal/repository"
	"stepik.leoscode.http/internal/service"
)

func (s *Server) CreateThread(w http.ResponseWriter, r *http.Request, params forum.CreateThreadParams) {
	if !validateXI(params.XIdempotencyKey) {
		writeError(w, http.StatusBadRequest, forum.ErrorUnauthorized{
			forum.ValidationError,
			&map[string]any{"XIdempotencyKey": params.XIdempotencyKey},
			"too long XIdempotencyKey"},
		)
		return
	}
	var threadCreate forum.ThreadCreate
	if err := readJson(r, &threadCreate); err != nil {
		writeError(w, http.StatusBadRequest, forum.ErrorBadRequest{
			forum.ValidationError,
			nil,
			"invalid json"},
		)
		return
	}
	if !validateThread(threadCreate) {
		writeError(w, http.StatusBadRequest, forum.ErrorBadRequest{
			forum.ValidationError,
			&map[string]any{"content": threadCreate.Content, "title": threadCreate.Title},
			"invalid params in thread"},
		)
		return
	}
	thread, err := s.service.CreateThread(threadCreate, repository.XUXI{
		XU: params.XUserId,
		XI: params.XIdempotencyKey,
	})
	switch {
	case err == nil:
		writeJson(w, http.StatusCreated, thread)
		return
	case errors.Is(err, service.ErrConflict):
		writeError(w, http.StatusConflict, forum.ErrorConflict{
			forum.BadRequest,
			nil,
			"distinct body but equal user and idempotency-key params in memory",
		})
	case errors.Is(err, service.ErrAlreadyThreadExist):
		writeJson(w, http.StatusOK, thread)
	case errors.Is(err, service.ErrNoSuchUserExist):
		writeError(w, http.StatusUnauthorized, forum.ErrorUnauthorized{
			forum.InvalidCredentials,
			nil,
			"no such user exist",
		})
	default:
		writeError(w, http.StatusInternalServerError, forum.ErrorInternal{
			forum.InternalError,
			nil,
			"no relevant error",
		})
	}

}

func (s *Server) GetThread(w http.ResponseWriter, r *http.Request, threadId forum.ThreadIdPath, params forum.GetThreadParams) {
	if !validateThreadId(threadId) {
		writeError(w, http.StatusBadRequest, forum.ErrorBadRequest{forum.ValidationError, &map[string]any{"thread id": threadId}, "thread is minimum 1"})
		return
	}
	t, err := s.service.GetThreadById(threadId, params.XUserId) // при постгресе может быть ошибка бд
	switch {
	case err == nil:
		writeJson(w, http.StatusOK, t)
	case errors.Is(err, service.ErrNoSuchUserExist):
		writeError(w, http.StatusUnauthorized, forum.ErrorUnauthorized{forum.Unauthorized, &map[string]any{}, "no thread found"})

	case errors.Is(err, service.ErrThreadNotFound):
		writeError(w, http.StatusNotFound, forum.ErrorNotFound{forum.NotFound, &map[string]any{"thread id": threadId}, "no thread found"})
	}
}

func validateXI(XI string) bool {
	return len(XI) >= 1 && len(XI) <= 128
}

func validateThread(thread forum.ThreadCreate) bool {
	isCorrectContent := len(thread.Content) >= 1 && len(thread.Content) <= 10000
	isCorrectTitle := len(thread.Title) >= 1 && len(thread.Title) <= 255
	return isCorrectContent && isCorrectTitle
}

func validateThreadId(id forum.ThreadIdPath) bool {
	return id >= 1
}
