package services

import (
	"backend/auth"
	"backend/database"
	"backend/models"
	"backend/models/utils"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
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

	user.Password, err = hashPassword(user.Password)

	if err != nil {
		log.Print("Failed to hash a password", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = database.SignUpStatement.QueryRow(user.Email, user.Password, user.FirstName, user.LastName).Scan(&user.Email)

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

	//Once the JWT pair is generated, we can store it using cookies
	accessCookie := &http.Cookie{
		Name:     "access-token",
		Value:    pair.Token,
		Expires:  time.Now().Add(time.Minute * 15),
		HttpOnly: true,
	}

	refreshCookie := &http.Cookie{
		Name:     "refresh-token",
		Value:    pair.RefreshToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}

	http.SetCookie(writer, accessCookie)
	http.SetCookie(writer, refreshCookie)

	writer.WriteHeader(http.StatusOK)

}

func RefreshToken(writer http.ResponseWriter, request *http.Request) {

	decoder := json.NewDecoder(request.Body)
	var jwtPair models.JWTPair
	var newPair models.JWTPair
	var response utils.GenericResponse
	var isValid = false

	err := decoder.Decode(&jwtPair)

	if err != nil {
		log.Print("Could not decode incoming refresh request", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	isValid, err = auth.ValidateToken(jwtPair.RefreshToken)

	if err != nil {
		log.Print("Failed jwt refresh request ", err)
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	if isValid {
		newPair, err = auth.GenerateJWTPair()

		if err != nil {
			log.Print("Failed to refresh JWT pair")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(newPair)
	} else {

		log.Print("Invalid refresh token received")
		response.Response = "Refresh token is incorrect"

		json.NewEncoder(writer).Encode(response)

		writer.WriteHeader(http.StatusUnauthorized)
		return

	}

}
