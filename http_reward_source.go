package mab

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func NewHTTPSource(client HttpDoer, url string, parser RewardParser, opts ...HTTPSourceOption) *HTTPSource {
	s := &HTTPSource{
		client:    client,
		url:       url,
		parser:    parser,
		marshaler: MarshalFunc(json.Marshal),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type HTTPSourceOption func(source *HTTPSource)

func WithContextMarshaler(m ContextMarshaler) HTTPSourceOption {
	return func(source *HTTPSource) {
		source.marshaler = m
	}
}

// HTTPSource is a basic implementation of RewardSource that gets reward estimates from an HTTP reward service.
type HTTPSource struct {
	client    HttpDoer
	url       string
	parser    RewardParser
	marshaler ContextMarshaler
}

func (h *HTTPSource) GetRewards(ctx context.Context, banditContext interface{}) ([]Dist, error) {
	body, err := h.marshaler.Marshal(banditContext)

	req, err := http.NewRequestWithContext(ctx, "POST", h.url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return h.parser.Parse(respBody)
}

type HttpDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type RewardParser interface {
	Parse([]byte) ([]Dist, error)
}

type ContextMarshaler interface {
	Marshal(v interface{}) ([]byte, error)
}

type ParseFunc func([]byte) ([]Dist, error)

func (p ParseFunc) Parse(b []byte) ([]Dist, error) { return p(b) }

type MarshalFunc func(v interface{}) ([]byte, error)

func (m MarshalFunc) Marshal(v interface{}) ([]byte, error) { return m(v) }

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
