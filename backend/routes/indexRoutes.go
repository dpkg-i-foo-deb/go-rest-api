package routes

import (
	"backend/app"
	"backend/services"
)

func InitIndexRoutes() {
	indexRoute()
}

func indexRoute() {
	app.AddGet("/", services.IndexService)
}
