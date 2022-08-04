package app

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

var App *fiber.App

func InitApp() {

	App = fiber.New()

}

func StartApp() {

	log.Fatal(App.Listen(os.Getenv("SERVER_PORT")))

}

func AddGet(route string, service func(connection *fiber.Ctx) error) {

	App.Get(route, service)

}
