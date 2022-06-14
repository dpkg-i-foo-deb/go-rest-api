package routes

import (
	"backend/auth"
	"backend/services"
)

func InitTaskRoutes() {
	createTaskRoute()
}

func createTaskRoute() {
	AddHandle("/create-task", auth.ValidateAndContinue(services.CreateTaskService))
}
