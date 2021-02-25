package mab

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProportional_ComputeProbs(t *testing.T) {
	tests := []struct {
		name     string
		rewards  []Dist
		expected []float64
	}{
		{
			"empty",
			[]Dist{},
			[]float64{},
		},
		{
			"one arm",
			[]Dist{Point(1)},
			[]float64{1},
		},
		{
			"null only",
			[]Dist{Null()},
			[]float64{0},
		},
		{
			"single point with nulls",
			[]Dist{Point(1), Null(), Null()},
			[]float64{1, 0, 0},
		},
		{
			"single normal",
			[]Dist{Normal(1, 5)},
			[]float64{1},
		},
		{
			"single normal",
			[]Dist{Normal(1, 5)},
			[]float64{1},
		},
		{
			"single beta",
			[]Dist{Beta(100, 150)},
			[]float64{1}},
		{
			"two points",
			[]Dist{
				Point(1),
				Point(3),
			},
			[]float64{.25, .75},
		},
		{
			"two points with nulls",
			[]Dist{
				Point(1),
				Null(),
				Null(),
				Point(3),
			},
			[]float64{.25, 0, 0, 0.75},
		},
		{
			"two normals",
			[]Dist{
				Normal(1, 1),
				Normal(3, 2),
			},
			[]float64{.25, 0.75},
		},
		{
			"two normals with nulls",
			[]Dist{
				Null(),
				Null(),
				Null(),
				Null(),
				Normal(1, 1),
				Normal(3, 2),
			},
			[]float64{0, 0, 0, 0, 0.25, 0.75},
		},
		{
			"two betas",
			[]Dist{
				Beta(10, 20), // mean = 1/3
				Beta(10, 5),  // mean = 2/3
			},
			[]float64{1. / 3, 2. / 3},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			strat := NewProportional()
			actual, err := strat.ComputeProbs(test.rewards)
			if err != nil {
				t.Fatal(err)
			}
			if !assert.ObjectsAreEqualValues(test.expected, actual) {
				t.Errorf("actual not %v, got=%v", test.expected, actual)
			}
		})
	}
}
