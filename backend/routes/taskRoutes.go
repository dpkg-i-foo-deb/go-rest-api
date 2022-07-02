package routes

import (
	"backend/auth"
	"backend/services"
)

func InitTaskRoutes() {
	createTaskRoute()
	getTaskRoute()
}

func createTaskRoute() {
	AddHandle("/tasks", auth.ValidateAndContinue(services.CreateTaskService), "PUT")
}

func getTaskRoute() {
	AddHandle("/tasks/{code}", auth.ValidateAndContinue(services.GetTaskService), "GET")
}
