package services

import (
	"backend/database"
	"backend/models"
	"backend/models/utils"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func SignUpService(writer http.ResponseWriter, request *http.Request) {

	decoder := json.NewDecoder(request.Body)
	var user models.User
	var response utils.SignUpResponse

	err := decoder.Decode(&user)

	if err != nil {
		log.Print("Could not decode incoming login request", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.Password, err = HashPassword(user.Password)

	if err != nil {
		log.Print("Failed to hash a password", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = database.SignUpStatement.QueryRow(user.Email, user.Password).Scan(&user.Email)

	if err != nil {
		log.Print("Failed to register a new user", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Email = user.Email
	response.Message = "Congrats! You can now log in!"

	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(response)

}
