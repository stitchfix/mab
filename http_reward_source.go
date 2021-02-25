package mab

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// NewHTTPSource returns a new HTTPSource given an HttpDoer, a url for the reward service, and a RewardParser.
// Optionally provide a ContextMarshaler for encoding bandit context.
// For example, if a reward service running on localhost:1337 provides Beta reward estimates:
//	client := &http.Client{timeout: time.Duration(100*time.Millisecond)}
//	url := "localhost:1337/rewards"
//	parser := ParseFunc(BetaFromJSON)
//	marshaler := MarshalFunc(json.Marshal)
//
//	source := NewHTTPSource(client, url, parser, WithContextMashaler(marshaler))
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

// HTTPSource is a basic implementation of RewardSource that gets reward estimates from an HTTP reward service.
type HTTPSource struct {
	client    HttpDoer
	url       string
	parser    RewardParser
	marshaler ContextMarshaler
}

// GetRewards makes a POST request to the reward URL, and parses the response into a []Dist.
// If a banditContext is provided, it will be marshaled and included in the body of the request.
func (h *HTTPSource) GetRewards(ctx context.Context, banditContext interface{}) ([]Dist, error) {

	var body io.Reader

	if banditContext != nil {
		marshaled, err := h.marshaler.Marshal(banditContext)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(marshaled)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", h.url, body)
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

// HTTPDoer is a basic interface for making HTTP requests. The net/http Client can be used or you can bring your own.
// Heimdall is a pretty good alternative client with some nice features: https://github.com/gojek/heimdall
type HttpDoer interface {
	Do(*http.Request) (*http.Response, error)
}

// RewardParser will be called to convert the response from the reward service to a slice of distributions.
type RewardParser interface {
	Parse([]byte) ([]Dist, error)
}

// ContextMarshaler is called on the banditContext and the result will become the body of the request to the bandit service.
type ContextMarshaler interface {
	Marshal(banditContext interface{}) ([]byte, error)
}

// HTTPSourceOption allows for optional arguments to NewHTTPSource
type HTTPSourceOption func(source *HTTPSource)

// WithContextMarshaler is an optional argument to HTTPSource
func WithContextMarshaler(m ContextMarshaler) HTTPSourceOption {
	return func(source *HTTPSource) {
		source.marshaler = m
	}
}

// ParseFunc is an adapter to allow a normal function to be used as a RewardParser
type ParseFunc func([]byte) ([]Dist, error)

func (p ParseFunc) Parse(b []byte) ([]Dist, error) { return p(b) }

// MarshalFunc is an adapter to allow a normal function to be used as a ContextMarshaler
type MarshalFunc func(banditContext interface{}) ([]byte, error)

func (m MarshalFunc) Marshal(banditContext interface{}) ([]byte, error) { return m(banditContext) }

// BetaFromJSON converts a JSON-encoded array of Beta distributions to a []Dist.
// Expects the JSON data to be in the form:
// 	`[{"alpha": 123, "beta": 456}, {"alpha": 3.1415, "beta": 9.999}]`
// Returns an error if alpha or beta value are missing or less than 1 for any arm.
// Any additional keys are ignored.
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

// NormalFromJSON converts a JSON-encoded array of Normal distributions to a []Dist.
// Expects the JSON data to be in the form:
// 	`[{"mu": 123, "sigma": 456}, {"mu": 3.1415, "sigma": 9.999}]`
// Returns an error if mu or sigma value are missing or sigma is less than 0 for any arm.
// Any additional keys are ignored.
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

// PointFromJSON converts a JSON-encoded array of Point distributions to a []Dist.
// Expects the JSON data to be in the form:
// 	`[{"mu": 123}, {"mu": 3.1415}]`
// Returns an error if mu value is missing for any arm. Any additional keys are ignored.
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
