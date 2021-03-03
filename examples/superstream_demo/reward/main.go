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
var campaignRewards map[int][]struct{ Alpha, Beta float64 }

func init() {
	campaignRewards = make(map[int][]struct{ Alpha, Beta float64 })

	campaignRewards[1] = []struct{ Alpha, Beta float64 }{
		{10, 125},
		{4, 130},
		{16, 80},
		{25, 99},
	}

	campaignRewards[2] = []struct{ Alpha, Beta float64 }{
		{25, 125},
		{5, 50},
		{7, 90},
		{13, 200},
	}
}

// This function handles incoming post requests to the /rewards endpoint
func handler(w http.ResponseWriter, r *http.Request) {
	
	// The request body must contain a JSON object with at least a "campaign_id" key and and integer value
	var req struct {
		CampaignID *int `json:"campaign_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.CampaignID == nil {
		http.Error(w, "missing required key \"campaign_id\"", http.StatusBadRequest)
		return
	}

	// get the context-dependent reward estimates
	rewards, ok := campaignRewards[*req.CampaignID]
	if !ok {
		http.Error(w, fmt.Sprintf("no rewards for campaign ID %d", req.CampaignID), http.StatusBadRequest)
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
