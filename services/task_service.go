package services

import (
	"encoding/base64"
	"errors"
	"fmt"
	"golang-todo-app-atp/repository"
	"golang-todo-app-atp/util"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type taskService struct {
	taskRepo repository.TaskRepository
}

func NewUserService(taskRepo repository.TaskRepository) TaskService {
	return taskService{taskRepo: taskRepo}
}

// Get
func (s taskService) GetTasks(search string, sort string) ([]TasksResponse, error) {
	tasks, err := s.taskRepo.GetAll(search, sort)
	if err != nil {
		logrus.Errorln("get tasks error: ", err)
		return nil, err
	}
	tasksRepo := []TasksResponse{}
	for _, task := range tasks {
		response := TasksResponse{
			Uid:         task.Uid,
			Title:       task.Title,
			Description: task.Description,
			CreatedAt:   task.CreatedAt,
			ImageUid:    task.ImageUid,
			Status:      task.Status,
		}
		tasksRepo = append(tasksRepo, response)
	}
	return tasksRepo, nil
}

func (s taskService) GetImageTask(uid uuid.UUID) (string, error) {
	imageTask, err := s.taskRepo.GetImage(&uid)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("Get ImageTask error %v", err))
		return "", err
	}

	imgByte, _ := base64.StdEncoding.DecodeString(imageTask.TaskImage.Image)
	if string(imgByte) != "" {
		return string(imgByte), nil
	}

	logrus.Errorln(fmt.Sprintf("Don't have image task"))
	return "", errors.New("Don't have ImageTask error")
}

// Create
func (s taskService) CreateTask(request NewTaskRequest) (*NewTaskRequestResponse, error) {
	data := request
	var imageUid *uuid.UUID

	validator := util.ValidateStruct(request)
	if validator != "" {
		logrus.Errorln(fmt.Sprintf("Body error %v", validator))
		return nil, errors.New(fmt.Sprintf("%v is required", validator))
	}

	// Add image
	if data.Image != "" && data.MimeType != "" {
		rawImage, err := util.CheckImage(data.Image, data.MimeType)
		if err != nil {
			logrus.Errorln("CheckImage Image error")
			return nil, err
		}

		if !util.CheckTaskStatus(data.Status) {
			logrus.Errorln("CheckTask Status error")
			return nil, errors.New("task status invalid")
		}

		createImage := repository.TaskImage{
			Image:    rawImage,
			MimeType: data.MimeType,
		}

		taskImages, err := s.taskRepo.CreateTaskImage(createImage)
		if err != nil {
			logrus.Errorln("Create images error")
			return nil, err
		}
		imageUid = &taskImages.Uid
	}
	// Create task
	createTask := repository.Task{
		Title:       data.Title,
		Description: data.Description,
		Status:      data.Status,
		ImageUid:    imageUid,
	}
	task, err := s.taskRepo.CreateTask(createTask)
	if err != nil {
		logrus.Errorln("Create Task error")
		return nil, err
	}
	response := NewTaskRequestResponse{
		Title:       task.Title,
		Status:      task.Status,
		Description: task.Description,
	}
	return &response, nil
}

// Update
func (s taskService) UpdateTask(uid uuid.UUID, request UpdateTaskRequest) error {
	data := request

	if !util.CheckTaskStatus(data.Status) {
		logrus.Errorln("invalid task status")
		return errors.New("task status invalid")
	}
	var taskBody repository.Task
	taskBody.Title = data.Title
	taskBody.Description = data.Description
	taskBody.Status = data.Status
	checkTaskImgUid, err := s.taskRepo.GetImage(&uid)
	if err != nil {
		logrus.Errorln("[update] get image error: ", err)
		return err
	}
	if data.Image != "" && data.MimeType != "" {
		rawImage, err := util.CheckImage(data.Image, data.MimeType)
		if err != nil {
			logrus.Errorln("image checking error: ", err)
			return err
		}
		if checkTaskImgUid.ImageUid == nil {
			createImage := repository.TaskImage{
				Image:    rawImage,
				MimeType: data.MimeType,
			}
			image, err := s.taskRepo.CreateTaskImage(createImage)
			if err != nil {
				logrus.Errorln("Create images error")
				return err
			}
			taskBody.ImageUid = &image.Uid
		} else { // Update old image
			var taskImage repository.TaskImage
			image, err := s.taskRepo.GetImage(checkTaskImgUid.ImageUid)
			if err != nil {
				logrus.Errorln("[update] get task error: ", err)
				return err
			}
			taskImage.Image = rawImage
			taskImage.MimeType = data.MimeType
			err = s.taskRepo.UpdateTaskImage(&image.Uid, taskImage)
			if err != nil {
				logrus.Errorln("update task image error: ", err)
				return err
			}
		}
	}
	if err := s.taskRepo.UpdateTask(&uid, taskBody); err != nil {
		logrus.Errorln("update task error: ", err)
		return err
	}

	return nil
}
