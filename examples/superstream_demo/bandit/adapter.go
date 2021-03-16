package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/stitchfix/mab"
)

type adapter struct {
	bandit mab.Bandit
}

type selectArmRequest struct {
	Unit    string          `json:"unit"`
	Context json.RawMessage `json:"context"`
}

type selectArmResponse struct {
	Rewards []mab.Dist `json:"rewards"`
	Probs   []float64  `json:"probs"`
	Arm     int        `json:"arm"`
}

func (a adapter) selectArm(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var req selectArmRequest

	if err := a.decodeRequestBody(r.Body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := a.bandit.SelectArm(r.Context(), req.Unit, req.Context)

	if err != nil {
		a.writeError(w, err)
		return
	}

	resp := selectArmResponse{
		Rewards: result.Rewards,
		Probs:   result.Probs,
		Arm:     result.Arm,
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func (a adapter) decodeRequestBody(b io.Reader, req *selectArmRequest) error {
	if err := json.NewDecoder(b).Decode(req); err == io.EOF {
		return fmt.Errorf("request body empty")
	} else if err != nil {
		return err
	}
	return nil
}

func (a adapter) writeError(w http.ResponseWriter, err error) {
	var non200 *mab.ErrRewardNon2XX
	if errors.As(err, &non200) {
		http.Error(w, err.Error(), err.(*mab.ErrRewardNon2XX).StatusCode)
		return
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
}
