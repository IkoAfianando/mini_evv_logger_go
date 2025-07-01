package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	_ "github.com/IkoAfianando/mini_evv_logger_go/docs"
	"github.com/IkoAfianando/mini_evv_logger_go/internal/router"
	"github.com/IkoAfianando/mini_evv_logger_go/internal/store"
)

// @title          Mini EVV Logger API
// @version        1.0
// @description    This is the API for the Caregiver Shift Tracker application.
// @contact.name   API Support
// @contact.email  support@bluehorntek.com
// @license.name   Apache 2.0
// @license.url    http://www.apache.org/licenses/LICENSE-2.0.html
// @host           localhost:8080
// @BasePath       /
// @schemes http
func main() {
	dataStore := store.NewStore()
	dataStore.SetupInitialData()
	app := fiber.New()
	router.SetupRoutes(app, dataStore)

	port := "8080"
	log.Printf("Starting server on port %s", port)
	log.Printf("Swagger UI is available at http://localhost:%s/swagger/index.html", port)
	log.Fatal(app.Listen(":" + port))
}
