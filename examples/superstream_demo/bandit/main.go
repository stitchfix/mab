package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/stitchfix/mab"
	"github.com/stitchfix/mab/numint"
)

func main() {
	cli := &http.Client{Timeout: time.Second}
	url := "http://reward-service/rewards"
	parser := mab.ParseFunc(mab.BetaFromJSON)
	marshaler := mab.MarshalFunc(json.Marshal)

	bandit := mab.Bandit{
		RewardSource: mab.NewHTTPSource(cli, url, parser, mab.WithContextMarshaler(marshaler)),
		Strategy:     mab.NewThompson(numint.NewQuadrature()),
		Sampler:      mab.NewSha1Sampler(),
	}

	r := mux.NewRouter()
	r.HandleFunc("/select_arm", adapter{bandit}.selectArm).Methods("POST")

	server := &http.Server{
		Handler: handlers.LoggingHandler(os.Stdout, r),
		Addr:    "0.0.0.0:80",
	}

	log.Fatal(server.ListenAndServe())
}
