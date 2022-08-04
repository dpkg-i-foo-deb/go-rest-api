package app

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

var App *fiber.App

func StartApp() {

	App = fiber.New()
	log.Fatal(App.Listen(":3000"))

}
