package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/stitchfix/mab"
	"github.com/stitchfix/mab/numint"
)

type randomizeRequest struct {
	Unit    string          `json:"unit"`
	Context json.RawMessage `json:"context"`
}

var bandit mab.Bandit

func handler(w http.ResponseWriter, r *http.Request) {
	var req randomizeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := bandit.SelectArm(r.Context(), req.Unit, req.Context)
	if err != nil {
		var non200 *mab.ErrRewardNon2XX
		if errors.As(err, &non200) {
			http.Error(w, err.Error(), err.(*mab.ErrRewardNon2XX).StatusCode)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(result)
}

func main() {
	cli := &http.Client{Timeout: time.Second}
	url := "http://reward-service/rewards"
	parser := mab.ParseFunc(mab.BetaFromJSON)
	marshaler := mab.MarshalFunc(json.Marshal)

	bandit = mab.Bandit{
		RewardSource: mab.NewHTTPSource(cli, url, parser, mab.WithContextMarshaler(marshaler)),
		Strategy:     mab.NewThompson(numint.NewQuadrature()),
		Sampler:      mab.NewSha1Sampler(),
	}

	r := mux.NewRouter()
	r.HandleFunc("/randomize", handler).Methods("POST")
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
