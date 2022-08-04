package services

import (
	"github.com/gofiber/fiber/v2"
)

func IndexService(connection *fiber.Ctx) error {

	return connection.SendString("Welcome to the API!")
}
