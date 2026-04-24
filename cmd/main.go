package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func webhookVerificationHandler(w http.ResponseWriter, r *http.Request) {
	challenge := r.URL.Query().Get("hub.challenge")

	if challenge == "" {
		http.Error(w, "missing challenge", http.StatusBadRequest)
		return
	}

	response := map[string]string{
		"hub.challenge": challenge,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /webhook/strava", webhookVerificationHandler)

	log.Println("Server running on :8080")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, mux)
}
