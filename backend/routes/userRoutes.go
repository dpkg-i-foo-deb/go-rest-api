package routes

import (
	"backend/auth"
	"backend/services"
)

func InitUserRoutes() {
	loginroute()
	signUpRoute()
	refreshRoute()
	signOutRoute()
}

func loginroute() {
	AddRoute("/login", services.LoginService, "POST", "OPTIONS")
}

func signUpRoute() {
	AddRoute("/sign-up", services.SignUpService, "POST", "OPTIONS")
}

func refreshRoute() {
	AddRoute("/refresh", services.RefreshToken, "GET", "OPTIONS")
}

func signOutRoute() {

	AddHandle("/sign-out", auth.ValidateAndContinue(services.SignOutService), "POST", "OPTIONS")

}
