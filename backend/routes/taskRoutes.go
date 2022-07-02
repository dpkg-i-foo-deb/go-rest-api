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
	AddHandle("/create-task", auth.ValidateAndContinue(services.CreateTaskService), "PUT")
}

func getTaskRoute() {
	AddHandle("/task/{code}", auth.ValidateAndContinue(services.GetTaskService), "GET")
}
