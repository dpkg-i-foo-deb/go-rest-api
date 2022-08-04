package app

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var app *fiber.App

func InitApp() {

	app = fiber.New()

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	app.Use(cors.New(cors.Config{
	
		AllowCredentials: true,
		AllowOrigins:"*",
		AllowMethods: "GET ,POST,OPTIONS,PUT,DELETE,PATCH",
		Next: nil,
		AllowHeaders: "",
		ExposeHeaders: "",
		MaxAge: 0,

	}))

}

func StartApp() {

	log.Fatal(app.Listen(os.Getenv("SERVER_PORT")))

}

func AddGet(route string, service func(connection *fiber.Ctx) error) {

	app.Get(route, service)

}

func AddPost(route string, service func(connection *fiber.Ctx) error) {

	app.Post(route, service)

}
