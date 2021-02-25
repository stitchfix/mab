package mab

import "testing"

func TestSha1Sampler_getIndex(t *testing.T) {
	tests := []struct {
		name     string
		weights  []float64
		bucket   int
		expected int
	}{
		{
			"single arm",
			[]float64{1.0},
			503,
			0,
		},
		{
			"two arms equal weight",
			[]float64{0.5, 0.5},
			0,
			0,
		},
		{
			"two arms equal weight",
			[]float64{0.5, 0.5},
			499,
			0,
		},
		{
			"two arms equal weight",
			[]float64{0.5, 0.5},
			500,
			1,
		},
		{
			"two arms equal weight",
			[]float64{0.5, 0.5},
			999,
			1,
		},
		{
			"three arms equal weight",
			[]float64{1.0, 1.0, 1.0},
			332,
			0,
		},
		{
			"three arms equal weight",
			[]float64{1.0, 1.0, 1.0},
			333,
			1,
		},
		{
			"three arms equal weight",
			[]float64{1.0, 1.0, 1.0},
			665,
			1,
		},
		{
			"three arms equal weight",
			[]float64{1.0, 1.0, 1.0},
			666,
			2,
		},
		{
			"three arms equal weight",
			[]float64{2.0, 2.0, 2.0},
			999,
			2,
		},
		{
			"zero weights",
			[]float64{0, 1, 0},
			0,
			1,
		},
		{
			"zero weights",
			[]float64{0, 1, 0},
			500,
			1,
		},
		{
			"zero weights",
			[]float64{0, 1, 0},
			999,
			1,
		},
		{
			"zero weights",
			[]float64{0, 1, 1},
			499,
			1,
		},
		{
			"zero weights",
			[]float64{0, 1, 1},
			500,
			2,
		},
	}

	s := NewSha1Sampler()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := s.getIndex(test.weights, test.bucket)
			if err != nil {
				t.Fatal(err)
			}
			if actual != test.expected {
				t.Errorf("index not %d, got=%d", test.expected, actual)
			}
		})
	}
}
