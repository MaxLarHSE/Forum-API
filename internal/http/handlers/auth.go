package handlers

import (
	"errors"
	"net/http"
	"regexp"

	forum "stepik.leoscode.http/internal/gen/api"
	"stepik.leoscode.http/internal/service"
)

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errResp := forum.ErrorResponse{
			Code:    forum.ValidationError,
			Details: nil,
			Message: err.Error(),
		}
		JsonEncode(w, errResp)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	if !ValidateUsername(username) {
		w.WriteHeader(http.StatusBadRequest)
		errResp := forum.ErrorBadRequest{
			Code:    forum.ValidationError,
			Details: &map[string]interface{}{"username": username},
			Message: "invalid username",
		}
		JsonEncode(w, errResp)
		return
	}
	if !ValidatePwd(password) {
		w.WriteHeader(http.StatusBadRequest)
		errResp := forum.ErrorBadRequest{
			Code:    forum.ValidationError,
			Details: &map[string]interface{}{"password": password},
			Message: "invalid password",
		}
		JsonEncode(w, errResp)
		return
	}
	id, err := s.service.Authorize(username, password)
	if errors.Is(err, service.InvalidPasswordError) {
		w.WriteHeader(http.StatusUnauthorized)
		errResp := forum.ErrorUnauthorized{
			Code:    forum.InvalidCredentials,
			Details: nil,
			Message: "incorrect password",
		}
		JsonEncode(w, errResp)
		return
	}
	w.WriteHeader(http.StatusOK)
	JsonEncode(w, map[string]any{"user_id": id})
}

var usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func ValidateUsername(username string) bool {
	return usernamePattern.MatchString(username) && len(username) >= 3 && len(username) <= 32
}

func ValidatePwd(pwd string) bool {
	return len(pwd) >= 8 && len(pwd) <= 64
}
