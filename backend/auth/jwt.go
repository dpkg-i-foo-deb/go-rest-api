package auth

import (
	"backend/models"
	"backend/models/utils"
	"backend/util"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type CustomClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWTPair(email string) (models.JWTPair, error) {

	var pair models.JWTPair

	log.Print("Generating new JWT pair")

	claims := CustomClaims{
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("AUTH_KEY")))
	if err != nil {
		log.Print("Could not generate a JWT pair")
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

		//If the requested method is options, the browser wants to negotiate CORS
		//We enable CORS to allow the frontend to make requests
		util.EnableCORS(&w)

		//If the requested method is options, the browser wants to negotiate CORS
		if r.Method == http.MethodOptions {
			//And we return 200 ok
			w.WriteHeader(http.StatusOK)
			return
		}

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

		accessCookie, err = r.Cookie("access-token")

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
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
