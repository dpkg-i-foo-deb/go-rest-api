package routes

import (
	"backend/app"
	"backend/auth"
	"backend/services"
)

func InitTaskRoutes() {
	createTaskRoute()
	getTaskRoute()
	getAllTasksRoute()
}

func createTaskRoute() {

	app.AddPut("/tasks", auth.ValidateAndContinue, services.CreateTaskService)

}

func getTaskRoute() {

	app.AddGet("/tasks/:code", auth.ValidateAndContinue, services.GetTaskService)

}

func getAllTasksRoute() {

	app.AddGet("/tasks", auth.ValidateAndContinue, services.GetAllTasksService)
}
