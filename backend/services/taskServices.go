package services

import (
	"backend/database"
	"backend/models"
	"encoding/json"
	"log"
	"net/http"
)

func CreateTask(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var task models.Task

	err := decoder.Decode(&task)

	if err != nil {
		log.Print("Could not decode incoming create task request ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = database.CreateTaskStatement.QueryRow(
		task.Title, task.Description, task.User, task.StartDate, task.DueDate, task.Status,
	).Scan(
		&task.Title, task.Code, task.Description, task.User, task.StartDate, task.DueDate, task.Status,
	)

	if err != nil {
		log.Print("Failed to create a new task ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(task)

}
