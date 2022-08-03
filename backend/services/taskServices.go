package services

import (
	"backend/auth"
	"backend/database"
	"backend/models"
	"backend/models/utils"
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateTaskService(writer http.ResponseWriter, request *http.Request, bodyBytes []byte) {

	var task *models.Task
	var tokenString string
	//Use the incoming request body bytes instead of the request which is already closed
	decoder := json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))

	err := decoder.Decode(&task)

	if err != nil {
		log.Print("Could not decode incoming create task request ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	tokenString, err = auth.GetCookieValue(request, "access-token")

	if err != nil {
		log.Print("Could not retrieve the token string ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	//We retrieve the token string from the request cookie
	task.UserEmail = new(string)
	*task.UserEmail, err = auth.EmailFromToken(tokenString)

	if err != nil {
		log.Print("Could not retrieve the user's identity ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Finally, we set the task email using the claims

	if err != nil {
		log.Print("Could not retrieve the auth cookie", err)
		writer.WriteHeader(http.StatusInternalServerError)
	}

	err = database.CreateTaskStatement.QueryRow(
		task.Title, task.Description, task.UserEmail, task.StartDate, task.DueDate, task.Status,
		task.MainTask,
	).Scan(
		&task.Title, &task.Description, &task.UserEmail,
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

func GetTaskService(writer http.ResponseWriter, request *http.Request, bodyBytes []byte) {

	taskCode, err := strconv.ParseInt(mux.Vars(request)["code"], 10, 64)

	var tokenString string

	var userEmail string

	var task models.Task

	var errorResponse utils.GenericResponse

	if err != nil {
		log.Print("The received code is not correct")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenString, err = auth.GetCookieValue(request, "access-token")

	if err != nil {
		log.Print("Could not retrieve the access token")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	userEmail, err = auth.EmailFromToken(tokenString)

	if err != nil {
		log.Print("The token does not contain the user email")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = database.GetTaskStatement.QueryRow(
		taskCode, userEmail,
	).Scan(
		&task.Title, &task.Description, &task.Code,
		&task.MainTask, &task.UserEmail, &task.StartDate, &task.DueDate, &task.Status,
	)

	if err != nil {
		log.Print("The queried task does not exist or the user has no access to it ", err)

		errorResponse.Response = "The task does not exist or you have no access to it"

		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode(errorResponse)

		return
	}

	writer.WriteHeader(http.StatusFound)
	json.NewEncoder(writer).Encode(task)

}

func GetAllTasksService(writer http.ResponseWriter, request *http.Request, bodyBytes []byte) {
	var tokenString string

	var userEmail string

	var task models.Task

	var allTasks []models.Task

	var errorResponse utils.GenericResponse

	var rows *sql.Rows

	tokenString, err := auth.GetCookieValue(request, "access-token")

	if err != nil {
		log.Print("Could not retrieve the access token")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	userEmail, err = auth.EmailFromToken(tokenString)

	if err != nil {
		log.Print("The token does not contain the user email")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	rows, err = database.GetAllTasksStatement.Query(userEmail)

	if err != nil {
		log.Print("Could not find any results for the request ", err)
		writer.WriteHeader(http.StatusNotFound)

		errorResponse.Response = "You have no tasks created"

		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode(errorResponse)

		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&task.Title, &task.Description, &task.Code,
			&task.MainTask, &task.UserEmail, &task.StartDate, &task.DueDate, &task.Status)

		if err != nil {
			log.Print("Could not retrieve the query result", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		allTasks = append(allTasks, task)
	}

	err = rows.Err()

	if err != nil {
		log.Print("Could not retrieve the query result", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(allTasks)
}

func EditTaskService(writer http.ResponseWriter, request *http.Request, bodyBytes []byte) {

	var tokenString string

	var userEmail string

	var task models.Task

	taskCode, err := strconv.ParseInt(mux.Vars(request)["code"], 10, 64)

	var errorResponse utils.GenericResponse

	var reader = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	if err != nil {
		log.Print("The received code is not valid")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenString, err = auth.GetCookieValue(request, "access-token")

	if err != nil {
		log.Print("Could not retrieve the access token")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	userEmail, err = auth.EmailFromToken(tokenString)

	if err != nil {
		log.Print("The token does not contain the user email")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(reader)

	err = decoder.Decode(&task)

	if err != nil {
		log.Print("Could not decode incoming request ", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	task.UserEmail = &userEmail

	code := int(taskCode)

	task.Code = &code

	err = database.EditTaskStatement.QueryRow(task.Title,
		task.Description, task.MainTask, task.StartDate,
		task.DueDate, task.Status, task.Code, task.UserEmail).Scan(
		&task.Title, &task.Description, &task.MainTask, &task.StartDate,
		&task.DueDate, &task.Status, &task.Code,
	)

	if err != nil {
		log.Print("Could not update a task ", err)
		writer.WriteHeader(http.StatusBadRequest)
		errorResponse.Response = "Could not update the task"

		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	writer.WriteHeader(http.StatusAccepted)
	json.NewEncoder(writer).Encode(task)

}

func DeleteTaskService(writer http.ResponseWriter, request *http.Request, bodyBytes []byte) {
	var tokenString string

	var userEmail string

	var task models.Task

	taskCode, err := strconv.ParseInt(mux.Vars(request)["code"], 10, 64)

	var errorResponse utils.GenericResponse

	if err != nil {
		log.Print("The received code is not valid")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenString, err = auth.GetCookieValue(request, "access-token")

	if err != nil {
		log.Print("Could not retrieve the access token")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	userEmail, err = auth.EmailFromToken(tokenString)

	if err != nil {
		log.Print("The token does not contain the user email")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	task.UserEmail = &userEmail

	code := int(taskCode)

	task.Code = &code

	result, err := database.DeleteTaskStatement.Exec(task.Code, task.UserEmail)

	if err != nil {
		log.Print("Could not delete a task ", err)
		writer.WriteHeader(http.StatusNotFound)

		errorResponse.Response = "Task not found or you have no access to it"

		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	affectedRows, err := result.RowsAffected()

	if err != nil || affectedRows == 0 {
		log.Print("Could not delete a task ", err)
		writer.WriteHeader(http.StatusNotFound)

		errorResponse.Response = "Task not found or you have no access to it"

		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	writer.WriteHeader(http.StatusOK)

}
