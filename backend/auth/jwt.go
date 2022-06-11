package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateJWT() (string, error) {
	log.Print("Creating a new Json Web Token")

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()

	tokenString, err := token.SignedString(os.Getenv("AUTH_KEY"))

	if err != nil {
		log.Print("Could not create a Json Web Token")
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("not authorized"))
				}
				return os.Getenv("AUTH_KEY"), nil
			})

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Not authorized: " + err.Error()))
			}

			if token.Valid {
				next(w, r)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Not authorized"))
		}
	})
}

func GetJWT(writer http.ResponseWriter, request *http.Request) {
	if request.Header["Api"] != nil {
		if request.Header["Api"][0] == os.Getenv("API_KEY") {
			token, err := CreateJWT()
			if err != nil {
				return
			}
			fmt.Fprint(writer, token)
		}
	}
}
