package handler

import (
	"errors"
	"golang-todo-app-atp/services"
	"golang-todo-app-atp/util"
	"strings"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type taskHandler struct {
	taskService services.TaskService
}

func NewTaskHandler(taskService services.TaskService) taskHandler {
	return taskHandler{taskService: taskService}
}

// Get
func (h taskHandler) GetTasks(ctx *fiber.Ctx) error {
	search := ctx.Query("search")
	sort := strings.ToLower(ctx.Query("sort"))
	tasks, err := h.taskService.GetTasks(search, sort)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(util.Result{
			Message: false,
			Error:   err.Error(),
		})
	}
	return ctx.JSON(util.Result{
		Message: true,
		Data:    tasks,
	})
}

func (h taskHandler) GetImageTask(ctx *fiber.Ctx) error {
	taskUid, err := uuid.FromString(ctx.Params("taskUid"))
	if err != nil {
		logrus.Errorln("invalid task uid")
		return ctx.Status(fiber.StatusBadRequest).JSON(util.Result{
			Message: false,
			Error:   errors.New("invalid task uid"),
		})
	}
	image, err := h.taskService.GetImageTask(taskUid)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(util.Result{
			Message: false,
			Error:   err.Error(),
		})
	}
	return ctx.SendString(image)
}

// Create
func (h taskHandler) CreateTask(ctx *fiber.Ctx) error {
	request := services.NewTaskRequest{}
	err := ctx.BodyParser(&request)
	if err != nil {
		logrus.Errorln("Body parser error")
		return ctx.Status(fiber.StatusBadRequest).JSON(util.Result{
			Message: false,
			Error:   err.Error(),
		})
	}

	task, err := h.taskService.CreateTask(request)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(util.Result{
			Message: false,
			Error:   err.Error(),
		})
	}
	return ctx.JSON(util.Result{
		Message: true,
		Data:    task,
	})
}

// Update
func (h taskHandler) Update(ctx *fiber.Ctx) error {
	taskUid, err := uuid.FromString(ctx.Params("taskUid"))
	if err != nil {
		logrus.Errorln("invalid task uid")
		return ctx.Status(fiber.StatusBadRequest).JSON(util.Result{
			Message: false,
			Error:   errors.New("invalid task uid"),
		})
	}

	request := services.UpdateTaskRequest{}
	if err := ctx.BodyParser(&request); err != nil {
		logrus.Errorln("body parser error: ", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(util.Result{
			Message: false,
			Error:   errors.New("body parser error"),
		})
	}

	err = h.taskService.UpdateTask(taskUid, request)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(util.Result{
			Message: false,
			Error:   err.Error(),
		})
	}

	return ctx.JSON(util.Result{
		Message:   true,
		MessageTh: "Update success",
	})
}
