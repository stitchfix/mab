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
		{
			"normals",
			[]Dist{
				Normal(1, 0.5),
				Normal(0.8, 0.44),
				Normal(2, 4.5),
				Normal(-1.5, 0.8),
				Normal(0, 0.8),
				Normal(4, 0.01),
			},
			[]float64{0, 0, 0.32832939702916675, 0, 0, 0.6715962238578759},
		},
		{
			"normals with nulls",
			[]Dist{
				Normal(1, 0.5),
				Normal(0.8, 0.44),
				Null(),
				Normal(2, 4.5),
				Normal(-1.5, 0.8),
				Normal(0, 0.8),
				Normal(4, 0.01),
				Null(),
			},
			[]float64{0, 0, 0, 0.32832939702916675, 0, 0, 0.6715962238578759, 0},
		},
		{
			"betas",
			[]Dist{
				Beta(100, 50),
				Beta(30, 100),
				Beta(5, 5),
				Beta(10, 5),
				Beta(20, 200),
			},
			[]float64{0.413633, 0, 0.098703, 0.487664, 0},
		},
		{
			"betas with null",
			[]Dist{
				Null(),
				Beta(100, 50),
				Beta(30, 100),
				Beta(5, 5),
				Beta(10, 5),
				Beta(20, 200),
			},
			[]float64{0, 0.413633, 0, 0.098703, 0.487664, 0},
		},
		{
			"lots of nulls",
			[]Dist{
				Null(),
				Null(),
				Null(),
				Null(),
				Null(),
				Null(),
				Null(),
				Null(),
				Null(),
				Null(),
				Null(),
				Null(),
				Null(),
				Null(),
				Null(),
				Beta(100, 100),
			},
			[]float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		},
		{
			"lots of nulls",
			[]Dist{
				Null(),
				Null(),
				Null(),
				Null(),
				Null(),
				Null(),
				Null(),
				Null(),
				Beta(30, 20),
				Null(),
				Null(),
				Null(),
				Null(),
				Null(),
				Null(),
				Beta(300, 300),
			},
			[]float64{0, 0, 0, 0, 0, 0, 0, 0, 0.915955668704697, 0, 0, 0, 0, 0, 0, 0.08442959550732976},
		},
		{
			"spicy betas",
			[]Dist{
				Beta(1988.9969421012, 21290.29165727936),
				Beta(50.513724206539536, 694.8915442828242),
				Beta(40.22907217881993, 474.05635888115313),
				Beta(63.51183105653544, 727.0899538364148),
				Beta(31.261111088044935, 411.1179082444311),
				Beta(21.92459706142498, 357.99764835992886),
				Beta(71.24351745432674, 818.4214002728952),
				Beta(52.28986733645648, 659.2207151426613),
				Beta(58.626012977120325, 718.5085688230059),
				Beta(27.76180147538136, 391.16613861489384),
			},
			[]float64{
				0.23448743303613864,
				0.015318543048527354,
				0.17017806247696898,
				0.17666880201042032,
				0.06656095618639102,
				0.008850309350875189,
				0.15737618680298987,
				0.058618352704498694,
				0.07651307073837736,
				0.035421658908035586,
			},
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

func BenchmarkThompson_ComputeProbs(b *testing.B) {
	rewards := []Dist{
		Beta(1988.9969421012, 21290.29165727936),
		Beta(50.513724206539536, 694.8915442828242),
		Beta(40.22907217881993, 474.05635888115313),
		Beta(63.51183105653544, 727.0899538364148),
		Beta(31.261111088044935, 411.1179082444311),
		Beta(21.92459706142498, 357.99764835992886),
		Beta(71.24351745432674, 818.4214002728952),
		Beta(52.28986733645648, 659.2207151426613),
		Beta(58.626012977120325, 718.5085688230059),
		Beta(27.76180147538136, 391.16613861489384),
	}
	startTol := 0.1
	endTol := 0.001
	for tol := startTol; tol >= endTol; tol /= 10 {
		bandit := NewThompson(numint.NewQuadrature(numint.WithAbsAndRelTol(tol, tol)))
		b.Run(fmt.Sprintf("tolerance_%v", tol), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := bandit.ComputeProbs(rewards)
				if err != nil {
					b.Error(err)
				}
			}
		})
	}

}

func closeEnough(a, b []float64) bool {
	tolerance := 0.001

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
