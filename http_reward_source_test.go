package mab

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBetaFromJSON(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected []Dist
	}{
		{
			"no arms",
			[]byte(`[]`),
			[]Dist{},
		},
		{
			"one arm",
			[]byte(`[{"alpha": 10, "beta": 20}]`),
			[]Dist{Beta(10, 20)},
		},
		{
			"lowercase",
			[]byte(`[{"alpha": 10, "beta": 20}, {"alpha": 20, "beta": 10}]`),
			[]Dist{Beta(10, 20), Beta(20, 10)},
		},
		{
			"mixed cases",
			[]byte(`[{"alpha": 10, "Beta": 20}, {"Alpha": 20, "beta": 10}]`),
			[]Dist{Beta(10, 20), Beta(20, 10)},
		},
		{
			"floats",
			[]byte(`[{"alpha": 10.0, "beta": 20.12345}, {"alpha": 1.945, "beta": 10}]`),
			[]Dist{Beta(10.0, 20.12345), Beta(1.945, 10)},
		},
		{
			"four arms",
			[]byte(`[{"alpha": 10.0, "beta": 20.12345}, {"alpha": 1.945, "beta": 10}, {"alpha": 100.0, "beta": 201.2345}, {"alpha": 999.9, "beta": 3.141}]`),
			[]Dist{Beta(10.0, 20.12345), Beta(1.945, 10), Beta(100.0, 201.2345), Beta(999.9, 3.141)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := BetaFromJSON(test.data)
			if err != nil {
				t.Fatal(err)
			}
			if !assert.ObjectsAreEqualValues(test.expected, actual) {
				t.Errorf("actual not %v. got=%v", test.expected, actual)
			}
		})
	}
}

func TestBetaFromJSONError(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			"empty response",
			[]byte(``),
		},
		{
			"not an array",
			[]byte(`{"alpha": 1, "beta": 2}`),
		},
		{
			"missing alpha",
			[]byte(`[{"alpha": 11.5, "beta": 25.0}, {"beta": 49.13}]`),
		},
		{
			"missing beta",
			[]byte(`[{"alpha": 11.5}, {"alpha": 11.5, "beta": 49.13}]`),
		},
		{
			"wrong params",
			[]byte(`[{"mu": 10, "sigma": 0.25}, {"mu": 0, "sigma": 0.8}]`),
		},
		{
			"alpha less than one",
			[]byte(`[{"alpha": -4, "beta": 20}, {"alpha": 200, "beta": 100}]`),
		},
		{
			"beta less than one",
			[]byte(`[{"alpha": 40, "beta": -0.1}, {"alpha": 200, "beta": 100}]`),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := BetaFromJSON(test.data)
			if err == nil {
				t.Error("expected error but didn't get one")
			}
		})
	}
}

func TestNormalFromJSON(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected []Dist
	}{
		{
			"no arms",
			[]byte(`[]`),
			[]Dist{},
		},
		{
			"one arm",
			[]byte(`[{"mu": 10, "sigma": 20}]`),
			[]Dist{Normal(10, 20)},
		},
		{
			"lowercase",
			[]byte(`[{"mu": 10, "sigma": 20}, {"mu": 20, "sigma": 10}]`),
			[]Dist{Normal(10, 20), Normal(20, 10)},
		},
		{
			"mixed cases",
			[]byte(`[{"mu": 10, "Sigma": 20}, {"Mu": 20, "sigma": 10}]`),
			[]Dist{Normal(10, 20), Normal(20, 10)},
		},
		{
			"floats",
			[]byte(`[{"mu": 10.0, "sigma": 20.12345}, {"mu": 1.945, "sigma": 10}]`),
			[]Dist{Normal(10.0, 20.12345), Normal(1.945, 10)},
		},
		{
			"negative mu",
			[]byte(`[{"mu": -10.0, "sigma": 20.12345}, {"mu": -1.945, "sigma": 10}]`),
			[]Dist{Normal(-10.0, 20.12345), Normal(-1.945, 10)},
		},
		{
			"four arms",
			[]byte(`[{"mu": 10.0, "sigma": 20.12345}, {"mu": 1.945, "sigma": 10}, {"mu": 100.0, "sigma": 201.2345}, {"mu": 999.9, "sigma": 3.141}]`),
			[]Dist{Normal(10.0, 20.12345), Normal(1.945, 10), Normal(100.0, 201.2345), Normal(999.9, 3.141)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := NormalFromJSON(test.data)
			if err != nil {
				t.Fatal(err)
			}
			if !assert.ObjectsAreEqualValues(test.expected, actual) {
				t.Errorf("actual not %v. got=%v", test.expected, actual)
			}
		})
	}
}

func TestNormalFromJSONError(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			"empty response",
			[]byte(``),
		},
		{
			"not an array",
			[]byte(`{"mu": 1, "sigma": 2}`),
		},
		{
			"missing mu",
			[]byte(`[{"mu": 11.5, "sigma": 25.0}, {"sigma": 49.13}]`),
		},
		{
			"missing sigma",
			[]byte(`[{"mu": 11.5}, {"mu": 11.5, "sigma": 49.13}]`),
		},
		{
			"wrong params",
			[]byte(`[{"alpha": 10, "beta": 0.25}, {"alpha": 0, "beta": 0.8}]`),
		},
		{
			"sigma less than one",
			[]byte(`[{"mu": -4, "sigma": 20}, {"mu": 200, "sigma": -100}]`),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NormalFromJSON(test.data)
			if err == nil {
				t.Error("expected error but didn't get one")
			}
		})
	}
}

func TestPointFromJSON(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected []Dist
	}{
		{
			"no arms",
			[]byte(`[]`),
			[]Dist{},
		},
		{
			"one arm",
			[]byte(`[{"mu": 10}]`),
			[]Dist{Point(10)},
		},
		{
			"lowercase",
			[]byte(`[{"mu": 10}, {"mu": 20}]`),
			[]Dist{Point(10), Point(20)},
		},
		{
			"mixed cases",
			[]byte(`[{"mu": 10}, {"Mu": 20}]`),
			[]Dist{Point(10), Point(20)},
		},
		{
			"floats",
			[]byte(`[{"mu": 10.0}, {"mu": 1.945}]`),
			[]Dist{Point(10.0), Point(1.945)},
		},
		{
			"negative mu",
			[]byte(`[{"mu": -10.0}, {"mu": -1.945, "sigma": 10}]`),
			[]Dist{Point(-10.0), Point(-1.945)},
		},
		{
			"four arms",
			[]byte(`[{"mu": 10.0}, {"mu": 1.945}, {"mu": 100.0}, {"mu": -999.9}]`),
			[]Dist{Point(10.0), Point(1.945), Point(100.0), Point(-999.9)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := PointFromJSON(test.data)
			if err != nil {
				t.Fatal(err)
			}
			if !assert.ObjectsAreEqualValues(test.expected, actual) {
				t.Errorf("actual not %v. got=%v", test.expected, actual)
			}
		})
	}
}

func TestPointFromJSONError(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			"empty response",
			[]byte(``),
		},
		{
			"not an array",
			[]byte(`{"mu": 1}`),
		},
		{
			"missing mu",
			[]byte(`[{"mu": 11.5}, {}]`),
		},
		{
			"wrong params",
			[]byte(`[{"alpha": 10, "beta": 0.25}, {"alpha": 0, "beta": 0.8}]`),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := PointFromJSON(test.data)
			if err == nil {
				t.Error("expected error but didn't get one")
			}
		})
	}
}
