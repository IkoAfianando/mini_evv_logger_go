package handler

import (
	"net/http"

	"github.com/IkoAfianando/mini_evv_logger_go/internal/router"
	"github.com/IkoAfianando/mini_evv_logger_go/internal/store"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	dataStore := store.NewStore()
	dataStore.SetupInitialData()
	app := fiber.New()
	router.SetupRoutes(app, dataStore)
	adaptor.FiberApp(app)(w, r)
}
