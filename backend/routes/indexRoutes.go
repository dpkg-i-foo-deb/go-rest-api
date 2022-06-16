package routes

import (
	"backend/auth"
	"backend/services"
)

func InitIndexRoutes() {
	indexRoute()
}

func indexRoute() {
	AddHandle("/", auth.ValidateAndContinue(services.IndexService), "GET")
}
