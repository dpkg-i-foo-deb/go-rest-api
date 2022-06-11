package services

import (
	"backend/auth"
	"backend/database"
	"backend/models"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
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

	//We recover both the user's email and password from database

	err = database.LoginStatement.QueryRow(user.Email).Scan(&queriedUser.Email, &queriedUser.Password)

	if err != nil {
		log.Print("Failed login attempt: ", err)
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte("username or password incorrect"))
		return
	}

	//We compare the stored password hash with its plain version

	err = bcrypt.CompareHashAndPassword([]byte(queriedUser.Password), []byte(user.Password))

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
