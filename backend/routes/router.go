package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var router *mux.Router

func InitRouter() {

	log.Print("Initializing router...")
	if router == nil {
		router = mux.NewRouter().StrictSlash(true)
	}
	log.Print("Router initialized!")
}

func GetRouter() *mux.Router {
	return router
}

func AddRoute(route string, function func(http.ResponseWriter, *http.Request)) {
	router.HandleFunc(route, function)
}
