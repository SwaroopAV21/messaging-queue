package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SwaroopAV21/messaging-queue/repository"
)

func getResponses(w http.ResponseWriter, r *http.Request) {
	responses, err := repository.GetResponses(r.Context())
	if err != nil {
		fmt.Println("retrieve responses:", err)
		http.Error(w, "could not retrieve responses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responses); err != nil {
		fmt.Println("encode responses:", err)
	}
}
