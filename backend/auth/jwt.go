package auth

import (
	"backend/models"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWTPair() (models.JWTPair, error) {

	var pair models.JWTPair

	log.Print("Generating new JWT pair")

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	t, err := token.SignedString([]byte(os.Getenv("AUTH_KEY")))
	if err != nil {
		log.Print("Could not generate a JWT pair")
		return pair, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	rt, err := refreshToken.SignedString([]byte(os.Getenv("AUTH_KEY")))
	if err != nil {
		log.Print("Could not generate a JWT pair")
		return pair, err
	}

	pair.RefreshToken = rt
	pair.Token = t

	return pair, nil
}
