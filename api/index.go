package api

import (
	"evv-logger-backend/internal/router"
	"evv-logger-backend/internal/store"
	"net/http"

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
