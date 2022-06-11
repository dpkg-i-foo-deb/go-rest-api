package services

import (
	"backend/auth"
	"backend/database"
	"backend/models"
	"encoding/json"
	"log"
	"net/http"
)

func LoginService(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var user models.User
	var queriedUser models.User
	var pair models.JWTPair

	err := decoder.Decode(&user)

	if err != nil {
		log.Print("Could not decode incoming login request", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = database.LoginStatement.QueryRow(user.Password, user.Email).Scan(&queriedUser.Email)

	if err != nil {
		log.Print("Failed login attempt: ", err)
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte("username or password incorrect"))
		return
	}

	pair, err = auth.GenerateJWTPair()

	if err != nil {
		log.Print("Login has failed_: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Internal server error"))
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(pair)

}
