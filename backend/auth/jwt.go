package auth

import (
	"backend/models"
	"backend/models/utils"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
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

func ValidateAndContinue(next func(writer http.ResponseWriter, request *http.Request, bodyBytes []byte)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//We gotta save the request body because you can only use it once
		bodyBytes, err := ioutil.ReadAll(r.Body)

		var tokenPair models.JWTPair
		var response utils.GenericResponse
		var isValid = false
		var accessCookie *http.Cookie

		//We must close the request body once we read it all
		r.Body.Close()

		//We also gotta check if we saved the body bytes correctly
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print("Could not save the request body bytes")
			response.Response = "Something went wrong, please try again"
			json.NewEncoder(w).Encode(response)
			return
		}

		//The decoder won't use the request body, it will use a new io stream
		decoder := json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))

		log.Print("Validating incoming request...")

		err = decoder.Decode(&tokenPair)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Print("The request contained invalid data")
			response.Response = "Invalid data was received"
			json.NewEncoder(w).Encode(response)
			return
		}

		accessCookie, err = r.Cookie("access-token")

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			log.Print("The request did not contain the access cookie")
			response.Response = "The access cookie was not found"
			json.NewEncoder(w).Encode(response)
			return
		}

		tokenPair.Token = accessCookie.Value

		isValid, err = ValidateToken(tokenPair.Token)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Print("The received token was invalid")
			response.Response = "Your token is invalid or has already expired"
			json.NewEncoder(w).Encode(response)
			return
		}

		if isValid {
			next(w, r, bodyBytes)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			log.Print("The received token was invalid")
			response.Response = "Your token is invalid or has already expired"
			json.NewEncoder(w).Encode(response)
			return
		}

	})
}
