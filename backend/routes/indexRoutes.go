package routes

import (
	"backend/auth"
	"backend/services"
)

func InitIndexRoutes() {
	indexRoute()
}

func indexRoute() {
	//AddRoute("/", services.IndexService)
	AddHandle("/", auth.ValidateJWT(services.IndexService))
}
