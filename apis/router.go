package apis

import (
	"fmt"
	"golang-todo-app-atp/handler"
	"golang-todo-app-atp/repository"
	"golang-todo-app-atp/services"
	"golang-todo-app-atp/util"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	getEnv = util.GoDotEnvVariable
	app    *fiber.App
)

func InitWebApi(db *gorm.DB) error {
	logrus.Infoln(fmt.Sprintf("Server start App %s Version %s at time: %s", getEnv("APP_NAME", ""), getEnv("VERSION", ""), time.Now().Format(time.RFC3339)))

	app = fiber.New()

	logOut := logrus.New()
	logOut.SetLevel(logrus.DebugLevel)
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "*",
		AllowMethods:     "*",
		AllowOrigins:     "*",
		AllowCredentials: true,
	}),
		logger.New(logger.Config{
			Format:     "${time} ${method} ${path} - ${ip} - ${status} - ${latency}\n",
			TimeZone:   "Asia/Bangkok",
			TimeFormat: time.ANSIC,
			Output:     logOut.WriterLevel(logrus.DebugLevel),
		}),
	)
	WebApiV1Endpoint(db)
	app.Listen(":8080")

	return nil
}

func WebApiV1Endpoint(db *gorm.DB) error {
	// repository
	taskRepository := repository.NewTaskRepository(db)
	// service
	taskService := services.NewUserService(taskRepository)
	// handler
	taskHandler := handler.NewTaskHandler(taskService)

	// Authen
	// app.Use("/api/users", jwtware.New(jwtware.Config{
	// 	SigningMethod:  "HS256",
	// 	SigningKey:     []byte(jwtSecret),
	// 	SuccessHandler: userHandler.AuthChk,
	// 	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
	// 		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"msg": false, "error": err.Error()})
	// 	},
	// }))

	api := app.Group("/api/v1")
	taks := api.Group("/task")
	taks.Get("", taskHandler.GetTasks)
	taks.Get("/image/:taskUid", taskHandler.GetImageTask)
	taks.Post("", taskHandler.CreateTask)
	taks.Put("/:taskUid", taskHandler.Update)
	return nil
}
