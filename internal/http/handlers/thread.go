package handlers

import (
	"errors"
	"net/http"

	forum "stepik.leoscode.http/internal/gen/api"
	"stepik.leoscode.http/internal/repository"
	"stepik.leoscode.http/internal/service"
)

func (s *Server) CreateThread(w http.ResponseWriter, r *http.Request, params forum.CreateThreadParams) {
	w.Header().Set("Content-Type", "application/json")
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

func validateXI(XI string) bool {
	return len(XI) >= 1 && len(XI) <= 128
}

func validateThread(thread forum.ThreadCreate) bool {
	isCorrectContent := len(thread.Content) >= 1 && len(thread.Content) <= 10000
	isCorrectTitle := len(thread.Title) >= 1 && len(thread.Title) <= 255
	return isCorrectContent && isCorrectTitle
}
