package services

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type TasksResponse struct {
	Uid         uuid.UUID  `json:"uid"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	ImageUid    *uuid.UUID `json:"image_uid"`
	Status      string     `json:"status"`
}

type TaskImageResponse struct {
	MimeType string `db:"mime_type"`
	Image    string `db:"image"`
}

type NewTaskRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	MimeType    string `json:"mime_type"`
	Status      string `json:"status" validate:"required"`
	Image       string `json:"image"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	MimeType    string `json:"mime_type"`
	Status      string `json:"status" validate:"required"`
	Image       string `json:"image"`
}

type NewTaskRequestResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type TaskService interface {
	// Get
	GetTasks(string, string) ([]TasksResponse, error)
	GetImageTask(uuid.UUID) (string, error)
	// Create
	CreateTask(NewTaskRequest) (*NewTaskRequestResponse, error)
	// Update
	UpdateTask(uuid.UUID, UpdateTaskRequest) error
}
