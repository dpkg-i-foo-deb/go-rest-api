package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

var router *mux.Router

func InitRouter() {

	if router == nil {
		router = mux.NewRouter().StrictSlash(true)
	}
}

func GetRouter() *mux.Router {
	return router
}

func AddRoute(route string, function func(http.ResponseWriter, *http.Request)) {
	router.HandleFunc(route, function)
}
