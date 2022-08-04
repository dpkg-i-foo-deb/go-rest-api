package services

import (
	"backend/auth"
	"backend/database"
	"backend/models"
	"backend/models/utils"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func SignUpService(connection *fiber.Ctx) error {

	decoder := json.NewDecoder(bytes.NewReader(connection.Body()))
	var user models.User
	var response utils.SignUpResponse

	err := decoder.Decode(&user)

	if err != nil {

		connection.Status(fiber.StatusBadRequest).SendString("Malformed request received")
		return errors.New("Bad signup request")

	}

	user.Password, err = hashPassword(user.Password)

	if err != nil {
		connection.Status(fiber.StatusInternalServerError).SendString("Failed to sign up, try again later")
		return errors.New("Failed to hash password")
	}

	err = database.SignUpStatement.QueryRow(user.Email, user.Password, user.FirstName, user.LastName).Scan(&user.Email)

	if err != nil {
		connection.Status(fiber.StatusInternalServerError).SendString("Failed to sign up, try again later")
		return errors.New("Failed to create user on database")
	}

	response.Email = user.Email
	response.Message = "Congrats! You can now log in!"

	connection.Status(fiber.StatusOK).JSON(response)

	return nil

}

func LoginService(connection *fiber.Ctx) error {

	decoder := json.NewDecoder(bytes.NewReader(connection.Body()))
	var user models.User
	var queriedUser models.User
	var pair models.JWTPair
	var response utils.GenericResponse

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

	response.Response = "Welcome!"

	connection.Status(fiber.StatusOK).JSON(response)
	return nil
}

func SignOutService(connection *fiber.Ctx) error {
	var response utils.GenericResponse

	var newRefreshCookie *fiber.Cookie
	var newAcessCookie *fiber.Cookie

	//Create the new cookies

	newAcessCookie = auth.GenerateFakeAccessCookie()
	newRefreshCookie = auth.GenerateFakeRefreshCookie()

	//Set the new cookies

	connection.Cookie(newAcessCookie)
	connection.Cookie(newRefreshCookie)

	response.Response = "Signed Out..."

	connection.Status(fiber.StatusOK).JSON(response)

	return nil

}

func RefreshTokenService(connection *fiber.Ctx) error {

	var response utils.GenericResponse
	refreshToken := connection.Cookies("refresh-token")
	var newPair models.JWTPair
	var userEmail string
	var newAccessCookie *fiber.Cookie
	var newRefreshCookie *fiber.Cookie

	response.Response = "The refresh token was not present"

	if refreshToken == "" {

		connection.Status(fiber.StatusUnauthorized).JSON(response)
		return nil
	}

	isValid, err := auth.ValidateToken(refreshToken)

	if isValid && err == nil {

		userEmail, err = auth.EmailFromToken(refreshToken)

		if err != nil {
			return errors.New("Failed to refresh token, try again later")
		}

		newPair, err = auth.GenerateJWTPair(userEmail)

		if err != nil {

			return errors.New("Failed to refresh token, try again later")
		}

		newAccessCookie = auth.GenerateAccessCookie(newPair.Token)
		newRefreshCookie = auth.GenerateRefreshCookie(newPair.RefreshToken)

		connection.Cookie(newAccessCookie)
		connection.Cookie(newRefreshCookie)

		return nil

	} else {
		response.Response = "The refresh token has expired or is not valid"
		connection.Status(fiber.StatusUnauthorized).JSON(response)
		return nil

	}

}
