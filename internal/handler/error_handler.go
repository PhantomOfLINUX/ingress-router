package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/PhantomOfLINUX/ingressRouter/internal/model"
)

func RespondWithError(w http.ResponseWriter, errorResponse model.ErrorResponse) {
    response, err := json.Marshal(errorResponse)
    if err != nil {
        log.Printf("Error marshaling JSON response: %v", err)
        RespondWithInternalServerError(w)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(errorResponse.StatusCode)
    w.Write(response)
}

func RespondWithInternalServerError(w http.ResponseWriter) {
    errorResponse := model.ErrorResponse{
        Response:   "error",
        Details:    "Internal Server Error",
        Error:      http.StatusText(http.StatusInternalServerError),
        StatusCode: http.StatusInternalServerError,
    }

    response, _ := json.Marshal(errorResponse)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusInternalServerError)
    w.Write(response)
}