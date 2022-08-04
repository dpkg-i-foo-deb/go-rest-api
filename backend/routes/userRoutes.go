package routes

import (
	"backend/app"
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

	app.AddPost("/login", services.LoginService)

}

func signUpRoute() {
	app.AddPost("/sign-up", services.SignUpService)
}

func signOutRoute() {

	app.AddPost("/sign-out", auth.ValidateAndContinue, services.SignOutService)
}

func refreshRoute() {

	app.AddGet("/refresh", services.RefreshTokenService)
}
