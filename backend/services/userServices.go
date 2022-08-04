package services

import (
	"backend/auth"
	"backend/database"
	"backend/models"
	"backend/models/utils"
	"backend/util"
	"bytes"
	"encoding/json"
	"errors"

	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
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

func LoginService(connection *fiber.Ctx) error {

	decoder := json.NewDecoder(bytes.NewReader(connection.Body()))
	var user models.User
	var queriedUser models.User
	var pair models.JWTPair

	err := decoder.Decode(&user)

	if err != nil {
		connection.Status(fiber.StatusInternalServerError).SendString("Error, try again later")
		return errors.New("Malformed Login Request Recieved")
	}

	//We recover both the user's email and password from database

	err = database.LoginStatement.QueryRow(user.Email).Scan(&queriedUser.Email, &queriedUser.Password)

	if err != nil {

		connection.Status(fiber.StatusUnauthorized).SendString("Username or password incorrect")
		return errors.New("Failed login attempt")
	}

	//We compare the stored password hash with its plain version

	err = bcrypt.CompareHashAndPassword([]byte(queriedUser.Password), []byte(user.Password))

	if err != nil {

		connection.Status(fiber.StatusUnauthorized).SendString("Username or password incorrect")
		return errors.New("Failed login attempt")
	}

	pair, err = auth.GenerateJWTPair(user.Email)

	if err != nil {

		connection.Status(fiber.StatusInternalServerError).SendString("Failed, try again later")
		return errors.New("Failed to generate JWT Pair")

	}

	//Once the JWT pair is generated, we can store it using cookies
	accessCookie := auth.GenerateAccessCookie(pair.Token)

	refreshCookie := auth.GenerateRefreshCookie(pair.RefreshToken)

	connection.Cookie(accessCookie)
	connection.Cookie(refreshCookie)

	connection.Status(fiber.StatusOK).SendString("Welcome!")
	return nil
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

}
