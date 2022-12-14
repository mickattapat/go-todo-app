package models

import (
	"golang-todo-app-atp/pkg/database"

	uuid "github.com/satori/go.uuid"
)

type Task struct {
	database.Model
	Title       string     `json:"title" gorm:"size:100"`
	Description *string    `json:"description" gorm:"type:TEXT"`
	Status      string     `json:"status"`
	ImageUid    *uuid.UUID `json:"image_uid" gorm:"size:191"`
	// preload
	TaskImage *TaskImage `json:"task_image" gorm:"foreignkey:image_uid;references:uid"`
}

type TaskImage struct {
	database.Model
	MimeType string `json:"mime_type" gorm:"type:TEXT"`
	Image    string `json:"image" gorm:"type:LONGTEXT"`
}
