package handlers

import (
	"encoding/json"
	"net/http"

	forum "stepik.leoscode.http/internal/gen/api"
	"stepik.leoscode.http/internal/service"
)

var _ forum.ServerInterface = (*Server)(nil)

type Server struct {
	service *service.Service
}

func NewServer(s *service.Service) *Server {
	return &Server{service: s}
}
func (s *Server) GetAttachment(w http.ResponseWriter, r *http.Request, attachmentId forum.AttachmentIdPath) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) DownloadAttachmentFile(w http.ResponseWriter, r *http.Request, attachmentId forum.AttachmentIdPath) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) Search(w http.ResponseWriter, r *http.Request, params forum.SearchParams) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) ListThreads(w http.ResponseWriter, r *http.Request, params forum.ListThreadsParams) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) DeleteThread(w http.ResponseWriter, r *http.Request, threadId forum.ThreadIdPath, params forum.DeleteThreadParams) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) PatchThread(w http.ResponseWriter, r *http.Request, threadId forum.ThreadIdPath, params forum.PatchThreadParams) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) ReplaceThread(w http.ResponseWriter, r *http.Request, threadId forum.ThreadIdPath, params forum.ReplaceThreadParams) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) UploadAttachment(w http.ResponseWriter, r *http.Request, threadId forum.ThreadIdPath, params forum.UploadAttachmentParams) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) ListPosts(w http.ResponseWriter, r *http.Request, threadId forum.ThreadIdPath, params forum.ListPostsParams) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) CreatePost(w http.ResponseWriter, r *http.Request, threadId forum.ThreadIdPath, params forum.CreatePostParams) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) HealthCheck(w http.ResponseWriter, r *http.Request, params forum.HealthCheckParams) {
	w.Header().Set("Content-Type", "application/json")
	if params.XRequestId != nil {
		w.Header().Set("X-Request-ID", *params.XRequestId)
	}
	if err := json.NewEncoder(w).Encode(map[string]any{"status": "ok"}); err != nil {
		http.Error(w, "MDA", http.StatusInternalServerError)
	}
}
