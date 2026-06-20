package handlers

import "net/http"

func (s *Server) InternalTruncate(w http.ResponseWriter, r *http.Request) {
	s.service.Truncate()

}
