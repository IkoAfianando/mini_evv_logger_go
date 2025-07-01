package router

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/IkoAfianando/mini_evv_logger_go/internal/handler"
	"github.com/IkoAfianando/mini_evv_logger_go/internal/store"
)

func SetupRoutes(app *fiber.App, st *store.Store) {
	scheduleHandler := handler.NewScheduleHandler(st)
	taskHandler := handler.NewTaskHandler(st)

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("EVV Logger Backend is running!")
	})
	api := app.Group("/api")

	// Admin route
	api.Post("/reset", scheduleHandler.ResetStore)

	// Schedule routes
	api.Get("/schedules", scheduleHandler.GetSchedules)
	api.Get("/schedules/today", scheduleHandler.GetTodaySchedules)
	api.Get("/schedules/:id", scheduleHandler.GetScheduleByID)

	// Visit routes
	api.Post("/schedules/:id/start", scheduleHandler.StartVisit)
	api.Post("/schedules/:id/end", scheduleHandler.EndVisit)
	api.Get("/schedules/:id/clock-in", scheduleHandler.ClockIn)
	api.Post("/schedules/:id/cancel-clock-in", scheduleHandler.CancelClockIn)

	// Task routes
	api.Post("/schedules/:id/tasks", scheduleHandler.AddTaskToSchedule)
	api.Put("/tasks/:taskId/update", taskHandler.UpdateTask)

	app.Get("/swagger/*", swagger.HandlerDefault)
}
