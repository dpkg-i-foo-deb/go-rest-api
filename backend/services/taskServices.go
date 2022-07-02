package services

import (
	"backend/auth"
	"backend/database"
	"backend/models"
	"backend/models/utils"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateTaskService(writer http.ResponseWriter, request *http.Request, bodyBytes []byte) {

	var task *models.Task
	//var claims *auth.CustomClaims
	var tokenString string
	//Use the incoming request body bites instead of the request which is already closed
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
	task.User = new(string)
	*task.User, err = auth.EmailFromToken(tokenString)

	if err != nil {
		log.Print("Could not retrieve the user's identity ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Finally, we set the task email using the claims

	//task.User = &claims.Email

	if err != nil {
		log.Print("Could not retrieve the auth cookie", err)
		writer.WriteHeader(http.StatusInternalServerError)
	}

	err = database.CreateTaskStatement.QueryRow(
		task.Title, task.Description, task.User, task.StartDate, task.DueDate, task.Status,
		task.MainTask,
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

func GetTaskService(writer http.ResponseWriter, request *http.Request, bodyByes []byte) {

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
		&task.MainTask, &task.User, &task.StartDate, &task.DueDate, &task.Status,
	)

	if err != nil {
		log.Print("The queried task does not exist or the user has no access to it ", err)

		errorResponse.Response = "The task does not exist or you have no access to it"

		json.NewEncoder(writer).Encode(errorResponse)

		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.WriteHeader(http.StatusFound)
	json.NewEncoder(writer).Encode(task)

}
