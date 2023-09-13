package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/nguyendinhduy02112002/user_todo/config"
	"github.com/nguyendinhduy02112002/user_todo/db"
	"github.com/nguyendinhduy02112002/user_todo/routes"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmfiber"
)

func setupRoutes(app *fiber.App) {
	// give response when at /
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "You are at the endpoint 😉",
		})
	})

	api := app.Group("/api")
	api.Get("", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "You are at the api endpoint 😉",
		})
	})

	routes.TodoRoute(api.Group("/user"))
}

func main() {
	tracer, err := apm.NewTracer("user_service", "2.0.0")
	if err != nil {
		fmt.Print(err)
	}
	defer tracer.Close()
	config.NewElasticsearchClient()
	config.ConnectDB()
	app := fiber.New()
	db.MigrateDB()
	app.Use(apmfiber.Middleware(apmfiber.WithTracer(tracer)))

	defer tracer.Flush(nil)
	setupRoutes(app)
	app.Listen(":3002")
}
