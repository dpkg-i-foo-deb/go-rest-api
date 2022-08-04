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
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,OPTIONS,PUT,DELETE,PATCH",
		Next:             nil,
		AllowHeaders:     "",
		ExposeHeaders:    "",
		MaxAge:           0,
	}))

}

func StartApp() {

	log.Fatal(app.Listen(os.Getenv("SERVER_PORT")))

}

func AddGet(route string, handlers ...fiber.Handler) {

	app.Get(route, handlers...)

}

func AddPost(route string, handlers ...fiber.Handler) {

	app.Post(route, handlers...)

}

func AddPut(route string, handlers ...fiber.Handler) {

	app.Put(route, handlers...)

}

func AddPatch(route string, handlers ...fiber.Handler){

	app.Patch(route,handlers...)
	
}

func AddDelete(route string, handlers ...fiber.Handler) {

	app.Delete(route, handlers...)

}
