package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
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
		return
	}

	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}

		log.Printf("Webhook received: %s", string(body))

		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			log.Printf("Invalid JSON: %v", err)
		} else {
			log.Printf("Parsed payload: %+v", payload)
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/webhook/strava", webhookHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
