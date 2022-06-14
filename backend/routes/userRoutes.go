package routes

import "backend/services"

func InitUserRoutes() {
	loginroute()
	signUpRoute()
	refreshRoute()
}

func loginroute() {
	AddRoute("/login", services.LoginService)
}

func signUpRoute() {
	AddRoute("/sign-up", services.SignUpService)
}

func refreshRoute() {
	AddRoute("/refresh", services.RefreshToken)
}
