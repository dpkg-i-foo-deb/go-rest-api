package routes

import (
	"backend/app"
	"backend/auth"
	"backend/services"
)

func InitTaskRoutes() {
	createTaskRoute()
}

func createTaskRoute() {

	app.AddPut("/tasks", auth.ValidateAndContinue, services.CreateTaskService)

}
