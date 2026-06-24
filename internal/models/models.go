package models

import "errors"

type ThreadPatchInput struct {
	Title    *string   `json:"title"`
	Content  *string   `json:"content"`
	Tags     *[]string `json:"tags"` // а зачем тут указатель интересно
	IsLocked *bool     `json:"is_locked"`
}

var (
	ErrNoJsonContentType = errors.New("No json content type")
)
