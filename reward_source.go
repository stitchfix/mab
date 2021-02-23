package mab

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HTTPRewardSource struct {
	client httpClient
	url    string
	parser rewardParser
}

func (h *HTTPRewardSource) GetRewards(context.Context) ([]Dist, error) {
	panic("not implemented")
}

type httpClient interface {
	Post(string, string, io.Reader) (*http.Response, error)
}

type rewardParser interface {
	Parse([]byte) ([]Dist, error)
}

type ParseFunc func([]byte) ([]Dist, error)

func (p ParseFunc) Parse(b []byte) ([]Dist, error) { return p(b) }

func BetaFromJSON(data []byte) ([]Dist, error) {
	var resp []struct {
		Alpha *float64 `json:"alpha"`
		Beta  *float64 `json:"beta"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	result := make([]Dist, len(resp))

	for i := range resp {
		if resp[i].Alpha == nil {
			return result, fmt.Errorf("missing alpha value for arm %d", i)
		}
		if resp[i].Beta == nil {
			return result, fmt.Errorf("missing beta value for arm %d", i)
		}
		if *resp[i].Alpha < 1 {
			return result, fmt.Errorf("arm %d alpha must be > 1. got=%f", i, *resp[i].Alpha)
		}
		if *resp[i].Beta < 1 {
			return result, fmt.Errorf("arm %d beta must be > 1. got=%f", i, *resp[i].Beta)
		}
		result[i] = Beta(*resp[i].Alpha, *resp[i].Beta)
	}

	return result, nil
}

func NormalFromJSON(data []byte) ([]Dist, error) {
	var resp []struct {
		Mu    *float64 `json:"mu"`
		Sigma *float64 `json:"sigma"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	result := make([]Dist, 0)

	for i := range resp {
		if resp[i].Mu == nil {
			return result, fmt.Errorf("missing mu value for arm %d", i)
		}
		if resp[i].Sigma == nil {
			return result, fmt.Errorf("missing sigma value for arm %d", i)
		}
		if *resp[i].Sigma < 0 {
			return result, fmt.Errorf("arm %d sigma must be > 0. got=%f", i, *resp[i].Sigma)
		}
		result = append(result, Normal(*resp[i].Mu, *resp[i].Sigma))
	}

	return result, nil
}

func PointFromJSON(data []byte) ([]Dist, error) {
	var resp []struct {
		Mu *float64
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	result := make([]Dist, 0)

	for i := range resp {
		if resp[i].Mu == nil {
			return result, fmt.Errorf("missing mu value for arm %d", i)
		}
		result = append(result, Point(*resp[i].Mu))
	}

	return result, nil
}
