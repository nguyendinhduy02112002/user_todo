package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm/module/apmhttp"
	"go.elastic.co/apm/v2"
	"golang.org/x/net/context/ctxhttp"

	"github.com/nguyendinhduy02112002/user_todo/config"
	"github.com/nguyendinhduy02112002/user_todo/models"
)

type Todo struct {
	ID        *int64    `json:"id"`
	UserID    *int64    `json:"user_id"`
	Title     *string   `json:"title"`
	Completed *bool     `json:"completed"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		apm.CaptureError(c.Context(), err).Send()
		return err
	}

	_, err := config.MI.DB.ExecContext(c.Context(), "INSERT INTO users (username, password) VALUES (?, ?)", user.Username, user.Password)
	if err != nil {
		apm.CaptureError(c.Context(), err).Send()
		return err
	}

	return c.SendString("User created successfully")
}

func GetUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	span, _ := apm.StartSpan(c.Context(), "processGetTodosByUserID", "custom")
	defer span.End()
	userInfo, err := CallTodosServiceAPI(userID, c)

	if err != nil {
		apm.CaptureError(c.Context(), err).Send()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error calling todos service"})
	}

	return c.JSON(userInfo)
}

func UpdateUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		apm.CaptureError(c.Context(), err).Send()
		return err
	}

	_, err := config.MI.DB.ExecContext(c.Context(), "UPDATE users SET username = ?, password = ? WHERE id = ?", user.Username, user.Password, userID)
	if err != nil {
		apm.CaptureError(c.Context(), err).Send()
		return err
	}

	return c.SendString("User updated successfully")
}

func DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	_, err := config.MI.DB.ExecContext(c.Context(), "DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		apm.CaptureError(c.Context(), err).Send()
		return err
	}

	return c.SendString("User deleted successfully")
}

func CallTodosServiceAPI(userID string, c *fiber.Ctx) ([]Todo, error) {

	url := fmt.Sprintf("http://localhost:3001/api/todos/user/%s", userID)
	client := apmhttp.WrapClient(&http.Client{})
	response, err := ctxhttp.Get(c.Context(), client, url)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		fmt.Println("Request failed with status code:", response.StatusCode)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	var todos []Todo
	err = json.Unmarshal(responseData, &todos)
	if err != nil {
		fmt.Println(err)
	}

	return todos, nil
}
