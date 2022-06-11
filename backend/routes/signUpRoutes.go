package routes

import "backend/services"

func InitSignUpRoutes() {
	signUpRoute()
}

func signUpRoute() {
	AddRoute("/sign-up", services.SignUpService)
}
