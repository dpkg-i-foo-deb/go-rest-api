package services

import (
	"backend/models/utils"
	"encoding/json"
	"log"
	"net/http"
)

func IndexService(writer http.ResponseWriter, request *http.Request, bodyBytes []byte) {

	var response utils.GenericResponse

	log.Printf("Index service accessed")

	response.Response = "success"

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)

}
