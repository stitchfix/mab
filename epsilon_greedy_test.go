package mab

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEpsilonGreedy_ComputeProbs(t *testing.T) {
	tests := []struct {
		name     string
		rewards  []Dist
		epsilon  float64
		expected []float64
	}{
		{
			"empty",
			[]Dist{},
			0.1,
			[]float64{},
		},
		{
			"single arm",
			[]Dist{Point(0.0)},
			0.1,
			[]float64{1},
		},
		{
			"null only",
			[]Dist{Null()},
			0.25,
			[]float64{0},
		},
		{
			"single point with nulls",
			[]Dist{Point(1), Null(), Null()},
			0.25,
			[]float64{1, 0, 0},
		},
		{
			"single normal",
			[]Dist{Normal(1, 5)},
			1,
			[]float64{1},
		},
		{
			"single normal",
			[]Dist{Normal(1, 5)},
			1,
			[]float64{1},
		},
		{
			"single beta",
			[]Dist{Beta(100, 150)},
			1,
			[]float64{1}},
		{
			"two points",
			[]Dist{
				Point(1),
				Point(3),
			},
			0.25,
			[]float64{.125, 0.875},
		},
		{
			"two points with nulls",
			[]Dist{
				Point(1),
				Null(),
				Null(),
				Point(3),
			},
			0.25,
			[]float64{.125, 0, 0, 0.875},
		},
		{
			"two normals",
			[]Dist{
				Normal(1, 1),
				Normal(3, 2),
			},
			0.25,
			[]float64{.125, 0.875},
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
			0.25,
			[]float64{0, 0, 0, 0, 0.125, 0.875},
		},
		{
			"two betas",
			[]Dist{
				Beta(10, 20), // mean = 1/3
				Beta(10, 5),  // mean = 3
			},
			0.25,
			[]float64{.125, 0.875},
		},
		{
			"negative reward",
			[]Dist{
				Point(-1),
				Point(3),
			},
			0.25,
			[]float64{.125, 0.875},
		},
		{
			"Epsilon zero",
			[]Dist{
				Point(-1),
				Point(3),
				Point(2.5),
			},
			0,
			[]float64{0, 1, 0},
		},
		{
			"Epsilon zero with null",
			[]Dist{
				Point(-1),
				Point(3),
				Null(),
				Point(2.5),
			},
			0,
			[]float64{0, 1, 0, 0},
		},
		{
			"Epsilon one",
			[]Dist{
				Point(-1),
				Point(3),
				Point(2.5),
			},
			1,
			[]float64{1.0 / 3, 1.0 / 3, 1.0 / 3},
		},
		{
			"Epsilon one with null",
			[]Dist{
				Point(-1),
				Point(3),
				Null(),
				Null(),
				Point(2.5),
				Null(),
			},
			1,
			[]float64{1.0 / 3, 1.0 / 3, 0, 0, 1.0 / 3, 0},
		},
		{
			"multiple maxima Epsilon zero",
			[]Dist{
				Point(-1),
				Point(3),
				Point(3),
			},
			0,
			[]float64{0, 0.5, 0.5},
		},
		{
			"multiple maxima Epsilon zero with null",
			[]Dist{
				Point(-1),
				Point(3),
				Null(),
				Point(3),
			},
			0,
			[]float64{0, 0.5, 0, 0.5},
		},
		{
			"all maxima Epsilon zero",
			[]Dist{
				Point(3),
				Point(3),
				Point(3),
				Point(3),
			},
			0,
			[]float64{0.25, 0.25, 0.25, 0.25},
		},
		{
			"all maxima Epsilon nonzero",
			[]Dist{
				Point(3),
				Point(3),
				Point(3),
				Point(3),
			},
			.5,
			[]float64{0.25, 0.25, 0.25, 0.25},
		},
		{"two maxima Epsilon nonzero",
			[]Dist{
				Point(3),
				Point(3),
				Point(1),
				Point(1),
			},
			.5,
			[]float64{0.375, 0.375, 0.125, 0.125},
		},
		{"two maxima Epsilon nonzero with null",
			[]Dist{
				Null(),
				Point(3),
				Point(3),
				Point(1),
				Point(1),
			},
			.5,
			[]float64{0, 0.375, 0.375, 0.125, 0.125},
		},
		{
			"three maxima Epsilon nonzero",
			[]Dist{
				Point(-4),
				Point(1),
				Point(1),
				Point(1),
			},
			.1,
			[]float64{0.025, 0.325, 0.325, 0.325},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			strat := NewEpsilonGreedy(test.epsilon)
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
