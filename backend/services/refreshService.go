package services

import (
	"backend/auth"
	"backend/models"
	"backend/models/utils"
	"encoding/json"
	"log"
	"net/http"
)

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
