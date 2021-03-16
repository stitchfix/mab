package main

import (
	"encoding/json"
	"errors"
	"fmt"
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
	r.HandleFunc("/select_arm", handler{bandit}.selectArm).Methods("POST")

	server := &http.Server{
		Handler: handlers.LoggingHandler(os.Stdout, r),
		Addr:    "0.0.0.0:80",
	}

	log.Fatal(server.ListenAndServe())
}

type handler struct {
	bandit mab.Bandit
}

type selectArmRequest struct {
	Unit    string          `json:"unit"`
	Context json.RawMessage `json:"context"`
}

func (h handler) selectArm(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var req selectArmRequest

	if err := h.decodeRequestBody(r.Body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.bandit.SelectArm(r.Context(), req.Unit, req.Context)

	if err != nil {
		h.writeError(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(result)
}

func (h handler) decodeRequestBody(b io.Reader, req *selectArmRequest) error {
	if err := json.NewDecoder(b).Decode(req); err == io.EOF {
		return fmt.Errorf("request body empty")
	} else if err != nil {
		return err
	}
	return nil
}

func (h handler) writeError(w http.ResponseWriter, err error) {
	var non200 *mab.ErrRewardNon2XX
	if errors.As(err, &non200) {
		http.Error(w, err.Error(), err.(*mab.ErrRewardNon2XX).StatusCode)
		return
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
}
