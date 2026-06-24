package handlers

import (
	"errors"
	"net/http"
	"regexp"

	forum "stepik.leoscode.http/internal/gen/api"
	"stepik.leoscode.http/internal/models"
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
	if err := readJson(r, &threadCreate); err != nil { // имеет ли смысл рассматривать ошибки другие
		writeError(w, http.StatusBadRequest, forum.ErrorBadRequest{
			forum.BadRequest,
			nil,
			"not json content type"},
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
	case errors.Is(err, service.ErrUserNotExist):
		writeError(w, http.StatusUnauthorized, forum.ErrorUnauthorized{
			forum.InvalidCredentials,
			nil,
			"no such user exist",
		})
	default:
		writeError(w, http.StatusInternalServerError, forum.ErrorInternal{
			forum.InternalError,
			&map[string]interface{}{"error": err},
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
	case errors.Is(err, service.ErrUserNotExist):
		writeError(w, http.StatusUnauthorized, forum.ErrorUnauthorized{forum.Unauthorized, &map[string]any{}, "no thread found"})

	case errors.Is(err, service.ErrThreadNotFound):
		writeError(w, http.StatusNotFound, forum.ErrorNotFound{forum.NotFound, &map[string]any{"thread id": threadId}, "no thread found"})
	default:
		writeError(w, http.StatusInternalServerError, forum.ErrorInternal{
			forum.InternalError,
			&map[string]interface{}{"error": err},
			"no relevant error",
		})
	}
}

func (s *Server) ListThreads(w http.ResponseWriter, r *http.Request, params forum.ListThreadsParams) {
	var limit, offset int32
	if params.Limit == nil {
		limit = 20
	} else {
		limit = *params.Limit
	}

	if params.Offset == nil {
		offset = 0
	} else {
		offset = *params.Offset
	}
	if !validateLimit(limit) {
		writeError(w, http.StatusBadRequest, forum.ErrorBadRequest{forum.ValidationError, &map[string]any{"limit": limit}, "limit must be >=1 and <=100"})
		return
	}
	if !validateOffset(offset) {
		writeError(w, http.StatusBadRequest, forum.ErrorBadRequest{forum.ValidationError, &map[string]any{"offset": offset}, "offset must be >=0"})
		return
	}
	if params.Tag != nil && !validateTag(*params.Tag) {
		writeError(w, http.StatusBadRequest, forum.ErrorBadRequest{forum.ValidationError, &map[string]any{}, "invalid tag"})
		return
	}
	//forum.ThreadListResponse{}
	threadListResp, err := s.service.GetListThreads(repository.ThreadListFilter{
		Limit:    limit,
		Offset:   offset,
		Tag:      params.Tag,
		AuthorID: params.AuthorId,
	})
	switch {
	case err == nil:
		writeJson(w, http.StatusOK, threadListResp)
	default:
		writeError(w, http.StatusInternalServerError, forum.ErrorInternal{
			forum.InternalError,
			&map[string]interface{}{"error": err},
			"no relevant error",
		})
	}
}

func (s *Server) ReplaceThread(w http.ResponseWriter, r *http.Request, threadId forum.ThreadIdPath, params forum.ReplaceThreadParams) {
	if !validateThreadId(threadId) {
		writeError(w, http.StatusBadRequest, forum.ErrorBadRequest{forum.ValidationError, &map[string]any{"thread id": threadId}, "thread is minimum 1"})
		return
	}
	var threadCreate forum.ThreadCreate
	if err := readJson(r, &threadCreate); err != nil {
		writeError(w, http.StatusBadRequest, forum.ErrorBadRequest{
			forum.BadRequest,
			nil,
			"not json content type"},
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

	thread, err := s.service.ReplaceThreadById(threadId, threadCreate, params)
	switch {
	case err == nil:
		writeJson(w, http.StatusOK, thread)
	case errors.Is(err, service.ErrUserDontHaveRights):
		writeError(w, http.StatusForbidden, forum.ErrorForbidden{forum.Forbidden, &map[string]any{}, "user dont have rights"})
	case errors.Is(err, service.ErrThreadNotFound):
		writeError(w, http.StatusNotFound, forum.ErrorNotFound{forum.NotFound, &map[string]any{"id": threadId}, "thread not found"})
	case errors.Is(err, service.ErrUserNotExist):
		writeError(w, http.StatusUnauthorized, forum.ErrorUnauthorized{forum.Unauthorized, &map[string]any{}, "user dont exist"})
	default:
		writeError(w, http.StatusInternalServerError, forum.ErrorInternal{
			forum.InternalError,
			&map[string]interface{}{"error": err}, // возможно не стоит ошибку передавать ибо она может содержать внутрянку
			"no relevant error",
		})
	}
}

func (s *Server) PatchThread(w http.ResponseWriter, r *http.Request, threadId forum.ThreadIdPath, params forum.PatchThreadParams) {
	if !validateThreadId(threadId) {
		writeError(w, http.StatusBadRequest, forum.ErrorBadRequest{forum.ValidationError, &map[string]any{"thread id": threadId}, "thread is minimum 1"})
		return
	}
	var threadPatch models.ThreadPatchInput
	if err := readJson(r, &threadPatch); err != nil {
		writeError(w, http.StatusBadRequest, forum.ErrorBadRequest{
			forum.BadRequest,
			nil,
			"not json content type"},
		)
		return
	}
	// в спеке вроде комменты о том что максимум одно поле может быть изменено за раз, хотя там anyOf...

	n := threadPatchFields(threadPatch)
	if n == 0 {
		writeError(w, http.StatusBadRequest, forum.ErrorBadRequest{forum.ValidationError, nil, "no fields to update"})
		return
	}

	if threadPatch.Title != nil {
		if !validateTitle(*threadPatch.Title) {
			writeError(w, http.StatusBadRequest, forum.ErrorBadRequest{forum.ValidationError, nil, "uncorrect title length"})
			return
		}
	}
	if threadPatch.Content != nil {
		if !validateContent(*threadPatch.Content) {
			writeError(w, http.StatusBadRequest, forum.ErrorBadRequest{forum.ValidationError, nil, "uncorrect content length"})
			return
		}
	}
	if threadPatch.Tags != nil && !validateTags(*threadPatch.Tags) {
		writeError(w, http.StatusBadRequest, forum.ErrorBadRequest{forum.ValidationError, nil, "uncorrect tags"})
		return
	}
	thread, err := s.service.ChangeThreadById(threadId, threadPatch, params)
	switch {
	case err == nil:
		writeJson(w, http.StatusOK, thread)
	case errors.Is(err, service.ErrUserDontHaveRights):
		writeError(w, http.StatusForbidden, forum.ErrorForbidden{forum.Forbidden, &map[string]any{}, "user dont have rights"})
	case errors.Is(err, service.ErrThreadNotFound):
		writeError(w, http.StatusNotFound, forum.ErrorNotFound{forum.NotFound, &map[string]any{"id": threadId}, "thread not found"})
	case errors.Is(err, service.ErrUserNotExist):
		writeError(w, http.StatusUnauthorized, forum.ErrorUnauthorized{forum.Unauthorized, &map[string]any{}, "user dont exist"})
	case errors.Is(err, service.ErrTryChangeLockedThread):
		writeError(w, http.StatusForbidden, forum.ErrorUnauthorized{forum.Forbidden, &map[string]any{}, "thread is locked"})
	default:
		writeError(w, http.StatusInternalServerError, forum.ErrorInternal{
			forum.InternalError,
			&map[string]interface{}{"error": err}, // возможно не стоит ошибку передавать ибо она может содержать внутрянку
			"no relevant error",
		})
	}
}

func (s *Server) DeleteThread(w http.ResponseWriter, r *http.Request, threadId forum.ThreadIdPath, params forum.DeleteThreadParams) {
	if !validateThreadId(threadId) {
		writeError(w, http.StatusBadRequest, forum.ErrorBadRequest{forum.ValidationError, &map[string]any{"thread id": threadId}, "thread is minimum 1"})
		return
	}

	err := s.service.DeleteThread(threadId, params) // тут именно парамс чтобы когда мы меняли параметры удаления то все было получше
	switch {
	case err == nil:
		writeJson(w, http.StatusNoContent, "тред удален")
	case errors.Is(err, service.ErrUserDontHaveRights): //404_delete_not_author почему требует 404 а не 403
		writeError(w, http.StatusNotFound, forum.ErrorNotFound{forum.NotFound, &map[string]any{}, "user dont have rights"})
	case errors.Is(err, service.ErrThreadNotFound):
		writeError(w, http.StatusNotFound, forum.ErrorNotFound{forum.NotFound, &map[string]any{"id": threadId}, "thread not found"})
	case errors.Is(err, service.ErrUserNotExist):
		writeError(w, http.StatusUnauthorized, forum.ErrorUnauthorized{forum.Unauthorized, &map[string]any{}, "user dont exist"})
	default:
		writeError(w, http.StatusInternalServerError, forum.ErrorInternal{
			forum.InternalError,
			&map[string]interface{}{"error": err}, // возможно не стоит ошибку передавать ибо она может содержать внутрянку
			"no relevant error",
		})
	}
}

func threadPatchFields(threadPatch models.ThreadPatchInput) int {
	n := 0
	if threadPatch.Title != nil {
		n++
	}
	if threadPatch.Content != nil {
		n++
	}
	if threadPatch.Tags != nil {
		n++
	}
	if threadPatch.IsLocked != nil {
		n++
	}
	return n
}

var tagRe = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

func validateTag(tag string) bool {
	return tagRe.MatchString(tag) && tag != "" && len(tag) <= 32
}

func validateTags(tags []string) bool {
	if tags == nil {
		return true
	}
	if len(tags) > 10 {
		return false
	}
	for i := range tags {
		if !validateTag(tags[i]) {
			return false
		}
	}
	return true
}
func validateLimit(l int32) bool {
	return l >= 1 && l <= 100
}
func validateOffset(o int32) bool {
	return o >= 0
}
func validateXI(XI string) bool {
	return len(XI) >= 1 && len(XI) <= 128
}
func validateThread(thread forum.ThreadCreate) bool {
	var validTags bool
	if thread.Tags == nil {
		validTags = true
	} else {
		validTags = validateTags(*thread.Tags)
	}
	return validateContent(thread.Content) && validateTitle(thread.Title) && validTags
}

var titleAndContentRe = regexp.MustCompile(`\s+`)

func validateContent(content string) bool {
	return len(content) >= 1 && len(content) <= 10000 && titleAndContentRe.ReplaceAllString(content, "") != ""
}
func validateTitle(title string) bool {
	return len(title) >= 1 && len(title) <= 255 && titleAndContentRe.ReplaceAllString(title, "") != ""
}

func validateThreadId(id forum.ThreadIdPath) bool {
	return id >= 1
}

// идея - вынести функцию универсальную для всех полей
