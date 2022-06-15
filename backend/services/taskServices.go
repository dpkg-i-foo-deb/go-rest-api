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
	body, err := ioutil.ReadAll(request.Body)

	if err != nil {
		log.Print("Could not save body contents ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	reader := bytes.NewReader(body)

	decoder := json.NewDecoder(reader)

	err = decoder.Decode(&task)

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
