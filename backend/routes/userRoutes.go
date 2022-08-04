package routes

import (
	"backend/app"
	"backend/services"
)

func InitUserRoutes() {
	loginroute()
	signUpRoute()
	//refreshRoute()
	signOutRoute()
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

	app.AddPost("/sign-out", services.SignOutService)
}
