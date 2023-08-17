package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"

	"github.com/nguyendinhduy02112002/user_todo/config"
	"github.com/nguyendinhduy02112002/user_todo/controller"
	"github.com/nguyendinhduy02112002/user_todo/db"

	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmfiber"
)

func main() {
	tracer, err := apm.NewTracer("users-service", "2.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer tracer.Close()

	config.NewElasticsearchClient()
	config.ConnectDB()
	app := fiber.New()
	db.MigrateDB()

	app.Use(apmfiber.Middleware(apmfiber.WithTracer(tracer)))

	setupUserAPI(app)
	defer tracer.Flush(nil)
	app.Listen(":3002")

}

func setupUserAPI(app *fiber.App) {
	app.Post("/api/users", controller.CreateUser)
	app.Get("/api/users/:id", controller.GetUser)
	app.Put("/api/users/:id", controller.UpdateUser)
	app.Delete("/api/users/:id", controller.DeleteUser)
}
