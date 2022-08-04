package app

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

var App *fiber.App

func StartApp() {

	App = fiber.New()
	log.Fatal(App.Listen(os.Getenv("PORT")))

}
