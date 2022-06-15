package services

import (
	"backend/database"
	"backend/models"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func CreateTaskService(writer http.ResponseWriter, request *http.Request, bodyBytes []byte) {

	var task models.Task

	//Use the incoming request body bites instead of the request which is already closed
	decoder := json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))

	err := decoder.Decode(&task)

	if err != nil {
		log.Print("Could not decode incoming create task request ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = database.CreateTaskStatement.QueryRow(
		task.Title, task.Description, task.User, task.StartDate, task.DueDate, task.Status,
	).Scan(
		&task.Title, &task.Description, &task.User,
		&task.StartDate, &task.DueDate, &task.Status, &task.MainTask, &task.Code,
	)

	if err != nil {
		log.Print("Failed to create a new task ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(task)

	log.Print("New task created!")

}
