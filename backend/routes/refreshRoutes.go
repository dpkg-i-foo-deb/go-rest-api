package routes

import "backend/services"

func InitRefreshRoutes() {
	refreshRoute()
}

func refreshRoute() {
	AddRoute("/refresh", services.RefreshToken)
}
