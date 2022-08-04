package main

import (
	"backend/database"
	"backend/routes"

	"log"

	"backend/app"

	"github.com/joho/godotenv"
)

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

}

func main() {

	initEnvironment()
	app.InitApp()
	database.InitDatabase()
	initQueries()
	initRoutes()

	app.StartApp()

}
