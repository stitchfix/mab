package mab

import (
	"testing"

	"github.com/stitchfix/mab"
	"github.com/stretchr/testify/assert"
)

func TestProportional_ComputeProbs(t *testing.T) {
	tests := []struct {
		name     string
		rewards  []mab.Dist
		expected []float64
	}{
		{
			"empty",
			[]mab.Dist{},
			[]float64{},
		},
		{
			"one arm",
			[]mab.Dist{mab.Point(1)},
			[]float64{1},
		},
		{
			"null only",
			[]mab.Dist{mab.Null()},
			[]float64{0},
		},
		{
			"single point with nulls",
			[]mab.Dist{mab.Point(1), mab.Null(), mab.Null()},
			[]float64{1, 0, 0},
		},
		{
			"single normal",
			[]mab.Dist{mab.Normal(1, 5)},
			[]float64{1},
		},
		{
			"single normal",
			[]mab.Dist{mab.Normal(1, 5)},
			[]float64{1},
		},
		{
			"single beta",
			[]mab.Dist{mab.Beta(100, 150)},
			[]float64{1}},
		{
			"two points",
			[]mab.Dist{
				mab.Point(1),
				mab.Point(3),
			},
			[]float64{.25, .75},
		},
		{
			"two points with nulls",
			[]mab.Dist{
				mab.Point(1),
				mab.Null(),
				mab.Null(),
				mab.Point(3),
			},
			[]float64{.25, 0, 0, 0.75},
		},
		{
			"two normals",
			[]mab.Dist{
				mab.Normal(1, 1),
				mab.Normal(3, 2),
			},
			[]float64{.25, 0.75},
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
			[]float64{0, 0, 0, 0, 0.25, 0.75},
		},
		{
			"two betas",
			[]mab.Dist{
				mab.Beta(10, 20), // mean = 1/3
				mab.Beta(10, 5),  // mean = 2/3
			},
			[]float64{1. / 3, 2. / 3},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			strat := mab.NewProportional()
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
