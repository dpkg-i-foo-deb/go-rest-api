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

func main() {

	log.Print("Initializing and loading environment")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Couldn't load .env file", err)
	}
	routes.InitRouter()
	database.InitDatabase()
	routes.InitIndexRoutes()

	var router = routes.GetRouter()

	fmt.Print("Server is running on port" + os.Getenv("SERVER_PORT") + "\n")
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_PORT"), router))

}
