package mab

import (
	"strconv"
	"testing"

	"github.com/stitchfix/mab"
	"gonum.org/v1/gonum/stat/distuv"
)

func TestSha1Sampler_Sample(t *testing.T) {

	tests := []struct {
		name    string
		weights []float64
	}{
		{
			"one weight",
			[]float64{1},
		},
		{
			"equal weights",
			[]float64{1.0, 1.0},
		},
		{
			"unequal weights",
			[]float64{2, 1},
		},
		{
			"odd number of weights",
			[]float64{1, 1, 1},
		},
		{
			"high and low probabilities",
			[]float64{0.001, 0.999},
		},
	}

	for _, test := range tests {

		var samples []int
		counts := make(map[int]float64)

		numIter := 10_000

		t.Run(test.name, func(t *testing.T) {
			s := mab.NewSha1Sampler()
			for i := 0; i < numIter; i++ {
				sample, err := s.Sample(test.weights, strconv.Itoa(i))
				if err != nil {
					t.Fatal(err)
				}
				samples = append(samples, sample)
				counts[sample] += 1
			}

			if len(test.weights) == 1 {
				testSingleWeight(t, samples)
				return
			}

			// for non-trivial examples we can do a statistical test on the frequencies
			testFrequencies(t, numIter, test.weights, counts)
		})
	}
}

func testSingleWeight(t *testing.T, samples []int) {
	for i, sample := range samples {
		if sample != 0 {
			t.Errorf("sample %d not 0. got=%v", i, sample)
		}
	}
}

func testFrequencies(t *testing.T, numIter int, weights []float64, observed map[int]float64) {
	if len(observed) != len(weights) {
		t.Fatalf("len(observed) != len(weights): %v, %v", observed, weights)
	}

	sumW := 0.0
	for _, w := range weights {
		sumW += w
	}

	expected := make(map[int]float64)
	for i, w := range weights {
		expected[i] = w * float64(numIter) / sumW
	}

	if len(observed) != len(expected) {
		t.Fatalf("len(observed) != len(expected): %v, %v", observed, expected)
	}

	chi2 := 0.0

	for val, obsCount := range observed {
		expCount, ok := expected[val]
		if !ok {
			t.Fatalf("missing expected fraction for value: %d. got=%v", val, expected)
		}

		chi2 += (obsCount - expCount) * (obsCount - expCount) / expCount
	}

	dof := float64(len(observed) - 1)
	dist := distuv.ChiSquared{K: dof}

	pVal := 1 - dist.CDF(chi2)
	alpha := 0.0001

	if pVal <= alpha {
		t.Errorf("expected frequencies %v, got=%v [pVal=%v]", expected, observed, pVal)
	}
}
