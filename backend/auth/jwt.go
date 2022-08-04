package auth

import (
	"backend/models"
	"backend/models/utils"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"log"
	"os"
	"time"
)

type CustomClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWTPair(email string) (models.JWTPair, error) {

	var pair models.JWTPair

	claims := CustomClaims{
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("AUTH_KEY")))
	if err != nil {

		return pair, err
	}

	claims = CustomClaims{
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	rt, err := refreshToken.SignedString([]byte(os.Getenv("AUTH_KEY")))
	if err != nil {

		return pair, err
	}

	pair.RefreshToken = rt
	pair.Token = t

	return pair, nil
}

func ValidateToken(tokenString string) (bool, error) {

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)

		if !ok {

			log.Print("Unexpected signing method")
			return nil, nil
		}

		return []byte(os.Getenv("AUTH_KEY")), nil
	})

	if err != nil {
		return false, errors.New("token is invalid")
	}

	if token.Valid {
		return true, err
	}

	return false, errors.New("token is not valid")
}

func ValidateAndContinue(connection *fiber.Ctx) error {

	accessToken := connection.Cookies("access-token")
	var response utils.GenericResponse
	response.Response = "The access token was not present"

	if accessToken == "" {
		connection.Status(fiber.StatusUnauthorized).JSON(response)
		return nil
	}

	isValid, err := ValidateToken(accessToken)

	if isValid && err == nil {
		connection.Next()
		return nil
	}
	response.Response = "The access token is not valid or has expired"
	connection.Status(fiber.StatusUnauthorized).JSON(response)
	return nil

}

func GetTokenClaims(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("AUTH_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	return token.Claims.(*CustomClaims), nil
}

func EmailFromToken(tokeString string) (string, error) {
	claims, err := GetTokenClaims(tokeString)

	if err != nil {
		return "", err
	}

	return claims.Email, nil
}
