package numint_test

import (
	"math"
	"testing"

	"github.com/stitchfix/mab"
	"github.com/stitchfix/mab/numint"
)

func TestQuadrature_Integrate(t *testing.T) {
	tests := []struct {
		f              func(float64) float64
		a, b, expected float64
	}{
		{
			func(x float64) float64 { return x },
			0, 0,
			0,
		},
		{
			func(x float64) float64 { return x },
			0, 1,
			0.5,
		},
		{
			func(x float64) float64 { return 1.0 / (1 + x*x) },
			0, 1,
			0.785398,
		},
		{
			mab.Beta(10, 20).Prob,
			0, 1,
			1,
		},
		{
			mab.Normal(10, 20).Prob,
			-700, 900,
			1,
		},
		{
			math.Asinh,
			-.5, 1,
			0.344588,
		},
		{
			func(x float64) float64 { return x * math.Cos(x*x) },
			1, 5,
			-0.486911,
		},
	}
	tol := 1E-6
	q := numint.NewQuadrature(numint.WithAbsTol(tol))

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			actual, err := q.Integrate(test.f, test.a, test.b)
			if err != nil {
				t.Fatal(err)
			}
			if !closeEnough(test.expected, actual, tol) {
				t.Errorf("actual not %f, got=%f", test.expected, actual)
			}
		})
	}
}

func closeEnough(a, b, tol float64) bool {
	return math.Abs(a-b) < tol
}
