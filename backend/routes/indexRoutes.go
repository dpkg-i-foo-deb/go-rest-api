package routes

import (
	"backend/auth"
	"backend/services"
)

func InitIndexRoutes() {
	indexRoute()
}

func indexRoute() {
	AddHandle("/", auth.ValidateJWT(services.IndexService))
}
