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

	database.InitLoginStatements()
	database.InitSignUpStatements()

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
	routes.InitLoginRoutes()
	routes.InitSignUpRoutes()
}

func startServer() {
	var router = routes.GetRouter()

	fmt.Print("Server is running on port" + os.Getenv("SERVER_PORT") + "\n")
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_PORT"), router))
}

func main() {

	initEnvironment()
	routes.InitRouter()
	database.InitDatabase()
	initQueries()
	initRoutes()
	startServer()

}
