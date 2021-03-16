package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// In a real system, rewards are stored in a DB, but for purposes of the demo we'll just hard-code some example values
var rewards map[int][]struct{ Alpha, Beta float64 }

func init() {
	rewards = make(map[int][]struct{ Alpha, Beta float64 })

	rewards[0] = []struct{ Alpha, Beta float64 }{
		{10, 125},
		{34, 130},
		{26, 95},
		{25, 99},
	}

	rewards[1] = []struct{ Alpha, Beta float64 }{
		{10, 125},
		{34, 130},
		{26, 95},
		{25, 99},
	}

	rewards[2] = []struct{ Alpha, Beta float64 }{
		{50, 250},
		{20, 105},
		{20, 75},
		{110, 399},
	}
}

// This function handles incoming post requests to the /rewards endpoint
func handler(w http.ResponseWriter, r *http.Request) {

	// The request body must contain a JSON object with at least a "source_id" key and and integer value
	var req struct {
		SourceID *int `json:"source_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err == io.EOF {
		http.Error(w, "request body empty", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.SourceID == nil {
		http.Error(w, "missing required key \"source_id\"", http.StatusBadRequest)
		return
	}

	// get the context-dependent reward estimates
	rewards, ok := rewards[*req.SourceID]
	if !ok {
		http.Error(w, fmt.Sprintf("no rewards for source ID %d", *req.SourceID), http.StatusBadRequest)
		return
	}

	// send a JSON-encoded response
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(rewards)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/rewards", handler).Methods("POST")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"alive": true}`)
	})

	server := &http.Server{
		Handler: handlers.LoggingHandler(os.Stdout, r),
		Addr:    "0.0.0.0:80",
	}

	log.Fatal(server.ListenAndServe())
}
