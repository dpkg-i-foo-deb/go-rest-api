package services

import (
	"backend/auth"
	"backend/database"
	"backend/models"
	"backend/models/utils"
	"backend/util"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func SignUpService(writer http.ResponseWriter, request *http.Request) {

	//We enable CORS to allow the frontend to make requests
	util.EnableCORS(&writer)

	//If the requested method is options, the browser wants to negotiate CORS
	if request.Method == http.MethodOptions {
		//And we return 200 ok
		writer.WriteHeader(http.StatusOK)
		return
	}

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

	//We enable CORS to allow the frontend to make requests
	util.EnableCORS(&writer)

	//If the requested method is options, the browser wants to negotiate CORS
	if request.Method == http.MethodOptions {
		//And we return 200 ok
		writer.WriteHeader(http.StatusOK)
		return
	}

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

	if err == sql.ErrNoRows {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

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

	pair, err = auth.GenerateJWTPair(user.Email)

	if err != nil {
		log.Print("Login has failed_: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Internal server error"))
		return
	}

	//Once the JWT pair is generated, we can store it using cookies
	accessCookie := auth.GenerateAccessCookie(pair.Token)

	refreshCookie := auth.GenerateRefreshCookie(pair.RefreshToken)

	http.SetCookie(writer, accessCookie)
	http.SetCookie(writer, refreshCookie)

	writer.WriteHeader(http.StatusOK)

}

func SignOutService(writer http.ResponseWriter, request *http.Request, bodyBytes []byte) {
	var response utils.GenericResponse

	var newRefreshCookie *http.Cookie
	var newAcessCookie *http.Cookie

	//Create the new cookies

	newAcessCookie = auth.GenerateFakeAccessCookie()
	newRefreshCookie = auth.GenerateFakeRefreshCookie()

	//Set the new cookies

	http.SetCookie(writer, newAcessCookie)
	http.SetCookie(writer, newRefreshCookie)

	writer.WriteHeader(http.StatusOK)

	response.Response = "Signed Out..."
	json.NewEncoder(writer).Encode(response)

	return

}

func RefreshToken(writer http.ResponseWriter, request *http.Request) {

	//We enable CORS to allow the frontend to make requests
	util.EnableCORS(&writer)

	//If the requested method is options, the browser wants to negotiate CORS
	if request.Method == http.MethodOptions {
		//And we return 200 ok
		writer.WriteHeader(http.StatusOK)
		return
	}

	var jwtPair models.JWTPair
	var newPair models.JWTPair
	var response utils.GenericResponse
	var isValid = false

	var newRefreshCookie *http.Cookie
	var newAcessCookie *http.Cookie
	var claims *auth.CustomClaims

	refreshCookie, err := request.Cookie("refresh-token")

	if err != nil {
		log.Print("The request did not contain a refresh cookie ", err)
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	jwtPair.RefreshToken = refreshCookie.Value

	isValid, err = auth.ValidateToken(jwtPair.RefreshToken)

	if err != nil {
		log.Print("Failed jwt refresh request ", err)
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	if isValid {

		//Retreieve the user email from the token
		claims, err = auth.GetTokenClaims(jwtPair.RefreshToken)

		if err != nil {
			log.Print("The token did not contain the required claims")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		newPair, err = auth.GenerateJWTPair(claims.Email)

		if err != nil {
			log.Print("Failed to refresh JWT pair")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		//Create the new cookies

		newAcessCookie = auth.GenerateAccessCookie(newPair.Token)
		newRefreshCookie = auth.GenerateRefreshCookie(newPair.RefreshToken)

		//Set the new cookies

		http.SetCookie(writer, newAcessCookie)
		http.SetCookie(writer, newRefreshCookie)

		writer.WriteHeader(http.StatusOK)

		response.Response = "refreshed"
		json.NewEncoder(writer).Encode(response)

	} else {

		log.Print("Invalid refresh token received")
		response.Response = "Refresh token is incorrect"

		json.NewEncoder(writer).Encode(response)

		writer.WriteHeader(http.StatusUnauthorized)
		return

	}

}
