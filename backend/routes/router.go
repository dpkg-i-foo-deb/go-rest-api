package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var Router *mux.Router

func InitRouter() {

	log.Print("Initializing router...")
	if Router == nil {
		Router = mux.NewRouter().StrictSlash(true)
	}
	log.Print("Router initialized!")
}

func GetRouter() *mux.Router {
	return Router
}

func AddRoute(route string, function func(http.ResponseWriter, *http.Request)) {
	Router.HandleFunc(route, function)
}

func AddHandle(route string, handle http.Handler) {
	Router.Handle(route, handle)
}
