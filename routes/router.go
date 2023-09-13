package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nguyendinhduy02112002/user_todo/controller"
)

func TodoRoute(route fiber.Router) {
	route.Put("/:id", controller.UpdateUser)
	route.Get("/:id", controller.GetUserDetail)
	route.Delete("/:id", controller.DeleteUser)
	route.Post("", controller.CreateUser)
	route.Get("", controller.GetAllUsers)
}
