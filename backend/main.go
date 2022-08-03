package main

import (
	"backend/database"
	"backend/routes"
	"fmt"
	"log"
	"net/http"
	"os"

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
	routes.InitTaskRoutes()
}

func startServer() {

	fmt.Print("Server is running on port" + os.Getenv("SERVER_PORT") + "\n")
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_PORT"), routes.Router))
}

func main() {

	initEnvironment()
	routes.InitRouter()
	database.InitDatabase()
	initQueries()
	initRoutes()
	startServer()

}
