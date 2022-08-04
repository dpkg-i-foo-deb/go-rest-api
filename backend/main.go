package main

import (
	"backend/database"
	"backend/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var App *fiber.App

func initQueries() {

	log.Print("Initializing database queries")

	database.InitUserStatements()
	database.InitTaskStatements()

	log.Print("Database queries initialized!")
}

func initEnvironment() {
	log.Print("Initializing and loading environment")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Couldn't load .env file", err)
	}
}

func initRoutes() {
	routes.InitIndexRoutes()
	routes.InitUserRoutes()
	routes.InitTaskRoutes()
}

func startServer() {
	App = fiber.New()
	log.Fatal(App.Listen(":3000"))
}

func main() {

	initEnvironment()
	//routes.InitRouter()
	database.InitDatabase()
	initQueries()
	//initRoutes()
	startServer()

}
