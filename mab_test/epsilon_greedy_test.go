package mab

import (
	"testing"

	"github.com/stitchfix/mab"
	"github.com/stretchr/testify/assert"
)

func TestEpsilonGreedy_ComputeProbs(t *testing.T) {
	tests := []struct {
		name     string
		rewards  []mab.Dist
		epsilon  float64
		expected []float64
	}{
		{
			"empty",
			[]mab.Dist{},
			0.1,
			[]float64{},
		},
		{
			"single arm",
			[]mab.Dist{mab.Point(0.0)},
			0.1,
			[]float64{1},
		},
		{
			"null only",
			[]mab.Dist{mab.Null()},
			0.25,
			[]float64{0},
		},
		{
			"single point with nulls",
			[]mab.Dist{mab.Point(1), mab.Null(), mab.Null()},
			0.25,
			[]float64{1, 0, 0},
		},
		{
			"single normal",
			[]mab.Dist{mab.Normal(1, 5)},
			1,
			[]float64{1},
		},
		{
			"single normal",
			[]mab.Dist{mab.Normal(1, 5)},
			1,
			[]float64{1},
		},
		{
			"single beta",
			[]mab.Dist{mab.Beta(100, 150)},
			1,
			[]float64{1}},
		{
			"two points",
			[]mab.Dist{
				mab.Point(1),
				mab.Point(3),
			},
			0.25,
			[]float64{.125, 0.875},
		},
		{
			"two points with nulls",
			[]mab.Dist{
				mab.Point(1),
				mab.Null(),
				mab.Null(),
				mab.Point(3),
			},
			0.25,
			[]float64{.125, 0, 0, 0.875},
		},
		{
			"two normals",
			[]mab.Dist{
				mab.Normal(1, 1),
				mab.Normal(3, 2),
			},
			0.25,
			[]float64{.125, 0.875},
		},
		{
			"two normals with nulls",
			[]mab.Dist{
				mab.Null(),
				mab.Null(),
				mab.Null(),
				mab.Null(),
				mab.Normal(1, 1),
				mab.Normal(3, 2),
			},
			0.25,
			[]float64{0, 0, 0, 0, 0.125, 0.875},
		},
		{
			"two betas",
			[]mab.Dist{
				mab.Beta(10, 20), // mean = 1/3
				mab.Beta(10, 5),  // mean = 3
			},
			0.25,
			[]float64{.125, 0.875},
		},
		{
			"negative reward",
			[]mab.Dist{
				mab.Point(-1),
				mab.Point(3),
			},
			0.25,
			[]float64{.125, 0.875},
		},
		{
			"Epsilon zero",
			[]mab.Dist{
				mab.Point(-1),
				mab.Point(3),
				mab.Point(2.5),
			},
			0,
			[]float64{0, 1, 0},
		},
		{
			"Epsilon zero with null",
			[]mab.Dist{
				mab.Point(-1),
				mab.Point(3),
				mab.Null(),
				mab.Point(2.5),
			},
			0,
			[]float64{0, 1, 0, 0},
		},
		{
			"Epsilon one",
			[]mab.Dist{
				mab.Point(-1),
				mab.Point(3),
				mab.Point(2.5),
			},
			1,
			[]float64{1.0 / 3, 1.0 / 3, 1.0 / 3},
		},
		{
			"Epsilon one with null",
			[]mab.Dist{
				mab.Point(-1),
				mab.Point(3),
				mab.Null(),
				mab.Null(),
				mab.Point(2.5),
				mab.Null(),
			},
			1,
			[]float64{1.0 / 3, 1.0 / 3, 0, 0, 1.0 / 3, 0},
		},
		{
			"multiple maxima Epsilon zero",
			[]mab.Dist{
				mab.Point(-1),
				mab.Point(3),
				mab.Point(3),
			},
			0,
			[]float64{0, 0.5, 0.5},
		},
		{
			"multiple maxima Epsilon zero with null",
			[]mab.Dist{
				mab.Point(-1),
				mab.Point(3),
				mab.Null(),
				mab.Point(3),
			},
			0,
			[]float64{0, 0.5, 0, 0.5},
		},
		{
			"all maxima Epsilon zero",
			[]mab.Dist{
				mab.Point(3),
				mab.Point(3),
				mab.Point(3),
				mab.Point(3),
			},
			0,
			[]float64{0.25, 0.25, 0.25, 0.25},
		},
		{
			"all maxima Epsilon nonzero",
			[]mab.Dist{
				mab.Point(3),
				mab.Point(3),
				mab.Point(3),
				mab.Point(3),
			},
			.5,
			[]float64{0.25, 0.25, 0.25, 0.25},
		},
		{"two maxima Epsilon nonzero",
			[]mab.Dist{
				mab.Point(3),
				mab.Point(3),
				mab.Point(1),
				mab.Point(1),
			},
			.5,
			[]float64{0.375, 0.375, 0.125, 0.125},
		},
		{"two maxima Epsilon nonzero with null",
			[]mab.Dist{
				mab.Null(),
				mab.Point(3),
				mab.Point(3),
				mab.Point(1),
				mab.Point(1),
			},
			.5,
			[]float64{0, 0.375, 0.375, 0.125, 0.125},
		},
		{
			"three maxima Epsilon nonzero",
			[]mab.Dist{
				mab.Point(-4),
				mab.Point(1),
				mab.Point(1),
				mab.Point(1),
			},
			.1,
			[]float64{0.025, 0.325, 0.325, 0.325},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			strat := mab.NewEpsilonGreedy(test.epsilon)
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
