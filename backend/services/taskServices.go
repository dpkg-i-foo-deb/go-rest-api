package services

import (
	"backend/auth"
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
	var claims *auth.CustomClaims
	var tokenString string
	//Use the incoming request body bites instead of the request which is already closed
	decoder := json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))

	err := decoder.Decode(&task)

	if err != nil {
		log.Print("Could not decode incoming create task request ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	//We retrieve the token string from the request cookie

	tokenString, err = auth.GetCookieValue(request, "access-token")

	if err != nil {
		log.Print("Could not retrieve the token string from the cookie ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Now we retrieve the claims from the token

	claims, err = auth.GetTokenClaims(tokenString)

	if err != nil {
		log.Print("Could not retrieve the claims from the token ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Finally, we set the task email using the claims

	task.User = &claims.Email

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
