package handlers

import (
	"net/http"

	forum "stepik.leoscode.http/internal/gen/api"
)

func (s *Server) InternalTruncate(w http.ResponseWriter, r *http.Request) {
	err := s.service.Truncate()
	if err != nil {
		writeError(w, http.StatusInternalServerError, forum.ErrorInternal{Code: forum.InternalError, Details: nil, Message: "invalid input"})
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
