package main

import (
	"backend/routes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Couldn't load .env file", err)
	}

	routes.InitRouter()
	routes.InitIndexRoutes()

	var router = routes.GetRouter()

	fmt.Print("Starting server... on port " + os.Getenv("SERVER_PORT") + "\n")
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_PORT"), router))

}
