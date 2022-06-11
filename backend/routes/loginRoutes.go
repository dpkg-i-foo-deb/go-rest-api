package routes

import "backend/services"

func InitLoginRoutes() {
	loginroute()
}

func loginroute() {
	AddRoute("/login", services.LoginService)
}
