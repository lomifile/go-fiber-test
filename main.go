package main

import (
	"github.com/auth-api/database"
	"github.com/auth-api/routes"
	"github.com/gofiber/fiber/v2"
  "github.com/gofiber/fiber/v2/middleware/cors"
  "github.com/swaggo/fiber-swagger"
  _"github.com/auth-api/docs"
)

// @title Auth Example API
// @version 1.0
// @description This is a sample Authentication app.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost
// @BasePath /api

func main() {
  database.Connect()
  app := fiber.New()
  app.Use(cors.New(cors.Config{
    AllowCredentials: true,
  }))
  app.Get("/docs/*", fiberSwagger.WrapHandler)
  routes.Setup(app)
  app.Listen(":8080")
}

