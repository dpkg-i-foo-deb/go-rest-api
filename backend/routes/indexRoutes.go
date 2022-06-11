package routes

import "backend/services"

func InitIndexRoutes() {
	indexRoute()
}

func indexRoute() {
	AddRoute("/", services.IndexService)
}
