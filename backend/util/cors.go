package util

import (
	"net/http"
	"os"
)

//Enable CORS and return the allowed hosts
func EnableCORS(writer *http.ResponseWriter) {

	//The URL should point to your frontend server
	//TODO move this URL to an environment variable
	(*writer).Header().Set("Access-Control-Allow-Origin", os.Getenv("ALLOWED_HOST"))
	(*writer).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	(*writer).Header().Set("Access-Control-Allow-Credentials", "true")
	(*writer).Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT,PATCH,DELETE")
}
