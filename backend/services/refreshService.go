package services

import (
	"backend/auth"
	"backend/models"
	"backend/models/utils"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

func RefreshToken(writer http.ResponseWriter, request *http.Request) {

	decoder := json.NewDecoder(request.Body)
	var jwtPair models.JWTPair
	var newPair models.JWTPair
	var response utils.GenericResponse

	err := decoder.Decode(&jwtPair)

	if err != nil {
		log.Print("Could not decode incoming refresh request", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := jwt.Parse(jwtPair.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			log.Print("Unexpected signing method detected")
			response.Response = "Unexpected signing method detected"

			json.NewEncoder(writer).Encode(response)
			writer.WriteHeader(http.StatusUnauthorized)

			return nil, nil
		}
		return []byte(os.Getenv("AUTH_KEY")), nil
	})

	if err != nil {
		log.Print("Failed jwt refresh request")
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	if token.Valid {
		newPair, err = auth.GenerateJWTPair()

		if err != nil {
			log.Print("Failed to refresh JWT pair")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(newPair)
	}

}
