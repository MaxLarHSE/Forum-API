package handlers

import (
	"errors"
	"net/http"
	"regexp"

	forum "stepik.leoscode.http/internal/gen/api"
	"stepik.leoscode.http/internal/service"
)

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		writeError(w, http.StatusBadRequest, forum.ErrorResponse{
			Code:    forum.ValidationError,
			Details: nil,
			Message: err.Error(),
		})
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	if !ValidateUsername(username) {
		writeError(w, http.StatusBadRequest, forum.ErrorBadRequest{
			Code:    forum.ValidationError,
			Details: &map[string]interface{}{"username": username},
			Message: "invalid username",
		})
		return
	}
	if !ValidatePwd(password) {
		writeError(w, http.StatusBadRequest, forum.ErrorBadRequest{
			Code:    forum.ValidationError,
			Details: &map[string]interface{}{"password": password},
			Message: "invalid password",
		})
		return
	}
	id, err := s.service.Authorize(username, password)
	if errors.Is(err, service.ErrInvalidPassword) {
		writeError(w, http.StatusUnauthorized, forum.ErrorUnauthorized{
			Code:    forum.InvalidCredentials,
			Details: nil,
			Message: "incorrect password",
		})
		return
	}
	writeJson(w, http.StatusOK, map[string]any{"user_id": id})
}

var usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func ValidateUsername(username string) bool {
	return usernamePattern.MatchString(username) && len(username) >= 3 && len(username) <= 32
}

func ValidatePwd(pwd string) bool {
	return len(pwd) >= 8 && len(pwd) <= 64
}
