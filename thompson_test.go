package mab

import (
	"fmt"
	"math"
	"testing"

	"github.com/stitchfix/mab/numint"
)

func ExampleThompson_ComputeProbs() {
	t := NewThompson(numint.NewQuadrature())
	rewards := []Dist{
		Beta(1989, 21290),
		Beta(40, 474),
		Beta(64, 730),
		Beta(71, 818),
		Beta(52, 659),
		Beta(59, 718),
	}
	probs, err := t.ComputeProbs(rewards)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.4f", probs)
	// Output: [0.2963 0.1760 0.2034 0.1690 0.0614 0.0939]
}

func TestThompson_ComputeProbs(t *testing.T) {
	tests := []struct {
		name     string
		rewards  []Dist
		expected []float64
	}{
		{
			"nil",
			nil,
			make([]float64, 0),
		},
		{
			"empty",
			make([]Dist, 0),
			make([]float64, 0),
		},
		{
			"single arm",
			[]Dist{Normal(0, 1.0)},
			[]float64{1},
		},
		{
			"equal arms",
			[]Dist{Normal(0, 1.0), Normal(0, 1.0)},
			[]float64{0.5, 0.5},
		},
		{
			"one null",
			[]Dist{Null()},
			[]float64{0},
		},
		{
			"several nulls",
			[]Dist{Null(), Null(), Null()},
			[]float64{0, 0, 0},
		},
		{
			"one non-null several nulls",
			[]Dist{Null(), Null(), Beta(10, 20), Null()},
			[]float64{0, 0, 1, 0},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := NewThompson(numint.NewQuadrature())
			actual, err := ts.ComputeProbs(test.rewards)
			if err != nil {
				t.Fatal(err)
			}
			if !closeEnough(test.expected, actual) {
				t.Errorf("actual not %v, got=%v", test.expected, actual)
			}
		})
	}
}

func closeEnough(a, b []float64) bool {
	tolerance := 0.0001

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		diff := math.Abs(a[i] - b[i])
		if diff > tolerance {
			return false
		}
	}
	return true
}
