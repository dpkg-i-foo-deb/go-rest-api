package routes

import (
	"backend/app"
	"backend/auth"
	"backend/services"
)

func InitUserRoutes() {
	loginroute()
	signUpRoute()
	//refreshRoute()
	//signOutRoute()
}

func loginroute() {

	app.AddPost("/login", services.LoginService)

}

func signUpRoute() {
	app.AddPost("/sign-up", services.SignUpService)
}

func refreshRoute() {
	AddRoute("/refresh", services.RefreshToken, "GET", "OPTIONS")
}

func signOutRoute() {

	AddHandle("/sign-out", auth.ValidateAndContinue(services.SignOutService), "POST", "OPTIONS")

}
