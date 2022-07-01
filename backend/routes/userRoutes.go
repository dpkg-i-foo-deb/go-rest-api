package routes

import "backend/services"

func InitUserRoutes() {
	loginroute()
	signUpRoute()
	refreshRoute()
}

func loginroute() {
	AddRoute("/login", services.LoginService, "POST")
}

func signUpRoute() {
	AddRoute("/sign-up", services.SignUpService, "POST")
}

func refreshRoute() {
	AddRoute("/refresh", services.RefreshToken, "GET")
}
