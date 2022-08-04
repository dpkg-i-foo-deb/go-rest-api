package routes

import (
	"backend/app"
	"backend/auth"
	"backend/services"
)

func InitTaskRoutes() {
	createTaskRoute()
	getTaskRoute()
}

func createTaskRoute() {

	app.AddPut("/tasks", auth.ValidateAndContinue, services.CreateTaskService)

}

func getTaskRoute(){

	app.AddGet("/tasks/:task-code",auth.ValidateAndContinue,services.GetTaskService)
	
}
