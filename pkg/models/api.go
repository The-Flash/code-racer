package models

type RuntimeResponse struct{}

type ExecutionRequest struct {
	Language   string          `json:"language" validate:"required"`
	EntryPoint string          `json:"entrypoint" validate:"required"`
	Files      []ExecutionFile `json:"files" validate:"required"`
}

type ExecutionFile struct {
	Name    string `json:"name" validate:"required"`
	Content string `json:"content" validate:"required"`
}
