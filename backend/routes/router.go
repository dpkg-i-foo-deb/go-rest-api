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

func AddRoute(route string, function func(http.ResponseWriter, *http.Request), methods ...string) {
	Router.HandleFunc(route, function).Methods(methods...)
}

func AddHandle(route string, handle http.Handler, methods ...string) {
	Router.Handle(route, handle).Methods(methods...)
}
