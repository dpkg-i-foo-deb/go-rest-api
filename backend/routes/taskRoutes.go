package routes

import (
	"backend/auth"
	"backend/services"
)

func InitTaskRoutes() {
	createTaskRoute()
	getTaskRoute()
	getAllTasksRoute()
	editTaskRoute()
	deleteTaskRoute()
}

func createTaskRoute() {
	AddHandle("/tasks", auth.ValidateAndContinue(services.CreateTaskService), "PUT", "OPTIONS")
}

func getTaskRoute() {
	AddHandle("/tasks/{code}", auth.ValidateAndContinue(services.GetTaskService), "GET", "OPTIONS")
}

func getAllTasksRoute() {
	AddHandle("/tasks", auth.ValidateAndContinue(services.GetAllTasksService), "GET", "OPTIONS")
}

func editTaskRoute() {
	AddHandle("/tasks/{code}", auth.ValidateAndContinue(services.EditTaskService), "PATCH", "OPTIONS")
}

func deleteTaskRoute() {
	AddHandle("/tasks/{code}", auth.ValidateAndContinue(services.DeleteTaskService), "DELETE", "OPTIONS")
}
