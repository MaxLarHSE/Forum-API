package models

type ThreadPatchInput struct {
	Title    *string   `json:"title"`
	Content  *string   `json:"content"`
	Tags     *[]string `json:"tags"`
	IsLocked *bool     `json:"is_locked"`
}
