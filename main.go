package main

import (
	"github.com/auth-api/database"
	"github.com/auth-api/routes"
	"github.com/gofiber/fiber/v2"
)


func main() {
  database.Connect()
  app := fiber.New()

  routes.Setup(app)
  app.Listen(":8080")
}

