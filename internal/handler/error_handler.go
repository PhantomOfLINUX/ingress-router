package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/PhantomOfLINUX/ingressRouter/internal/model"
)

func RespondWithError(w http.ResponseWriter, statusCode int, response, errorCode, details string) {
	errorResponse := model.ErrorResponse{
		Response: response,
		Error:    errorCode,
		Details:  details,
	}

	jsonResponse, err := json.Marshal(errorResponse)
	if err != nil {
		log.Printf("Error marshaling JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}