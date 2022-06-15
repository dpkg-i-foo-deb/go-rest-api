package services

import (
	"fmt"
	"log"
	"net/http"
)

func IndexService(writer http.ResponseWriter, request *http.Request, bodyBytes []byte) {
	log.Printf("Index service accessed")
	fmt.Fprint(writer, "Welcome!")
}
