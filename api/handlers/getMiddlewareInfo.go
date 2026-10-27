package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/SwaroopAV21/messaging-queue/repository"
)


func  getResponses(w http.ResponseWriter, r *http.Request) {
	responses, err := repository.GetResponses(r.Context())
    if err != nil {
        http.Error(w, "could not retrieve responses", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(responses); err != nil {
        log.Printf("encode responses: %v", err)
    }
}

func getRequests(w http.ResponseWriter, r *http.Request) {
	responses, err := repository.GetRequests(r.Context())
    if err != nil {
		fmt.Print("error: ", err)
        http.Error(w, "could not retrieve responses", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(responses); err != nil {
        log.Printf("encode responses: %v", err)
    }
}