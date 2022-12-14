package repository

import (
	"golang-todo-app-atp/pkg/database"

	uuid "github.com/satori/go.uuid"
)

type Task struct {
	database.Model
	Title       string     `db:"title"`
	Description string     `db:"description" jsonschema:"required"`
	ImageUid    *uuid.UUID `db:"image_uid"`
	Status      string     `db:"status"`
	// preload
	TaskImage TaskImage `db:"task_image" gorm:"foreignkey:image_uid;references:uid"`
}

type TaskImage struct {
	database.Model
	MimeType string `db:"mime_type"`
	Image    string `db:"image"`
}

type TaskRepository interface {
	// Get
	GetAll(string, string) ([]Task, error)
	GetImage(*uuid.UUID) (*Task, error)
	// Create
	CreateTask(Task) (*Task, error)
	CreateTaskImage(TaskImage) (*TaskImage, error)
	// Update
	UpdateTask(*uuid.UUID, Task) error
	UpdateTaskImage(*uuid.UUID, TaskImage) error
}
