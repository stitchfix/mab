package numint

func NewtonCotesOpen(degree int) NewtonCotesRule {
	var coeffs []float64
	switch degree {
	default:
	case 2:
		coeffs = []float64{2.0}
	case 3:
		coeffs = []float64{3. / 2, 3. / 2}
	case 4:
		coeffs = []float64{8. / 3, -4. / 3, 8. / 3}
	case 5:
		coeffs = []float64{55. / 24, 5. / 24, 5. / 24, 55. / 24}
	case 6:
		coeffs = []float64{72. / 35, 27. / 35, 12. / 35, 27. / 35, 72. / 35}
	case 7:
		coeffs = []float64{91. / 48, 49. / 48, 28. / 48, 42. / 48, 49. / 48, 91. / 48}
	}
	return NewtonCotesRule{coeffs: coeffs, open: true}
}

func NewtonCotesClosed(degree int) NewtonCotesRule {
	var coeffs []float64
	switch degree {
	default:
	case 1:
		coeffs = []float64{0.5, 0.5}
	case 2:
		coeffs = []float64{1. / 3, 4. / 3, 1. / 3}
	case 3:
		coeffs = []float64{3. / 8, 9. / 8, 9. / 8, 3. / 8}
	case 4:
		coeffs = []float64{14. / 45, 64. / 45, 24. / 45, 64. / 45, 14. / 45}
	case 5:
		coeffs = []float64{95. / 288, 375. / 288, 250. / 288, 250. / 288, 375. / 288, 95. / 288}
	}
	return NewtonCotesRule{coeffs: coeffs, open: false}
}

type NewtonCotesRule struct {
	coeffs []float64
	open   bool
}

func (n *NewtonCotesRule) Weights(a float64, b float64) []float64 {
	weights := make([]float64, len(n.coeffs))
	for i := range n.coeffs {
		weights[i] = n.stepSize(a, b) * n.coeffs[i]
	}
	return weights
}

func (n *NewtonCotesRule) Points(a float64, b float64) []float64 {
	if n.degree() <= 0 {
		return []float64{}
	}

	if n.open {
		return n.openPoints(a, b)
	}

	return n.closedPoints(a, b)
}

func (n *NewtonCotesRule) openPoints(a float64, b float64) []float64 {
	result := make([]float64, len(n.coeffs))
	h := n.stepSize(a, b)

	x := a + h

	for i := range result {
		result[i] = x
		x += h
	}
	return result
}

func (n *NewtonCotesRule) closedPoints(a float64, b float64) []float64 {
	result := make([]float64, len(n.coeffs))
	h := n.stepSize(a, b)

	x := a

	for i := range result {
		result[i] = x
		x += h
	}
	return result
}

func (n *NewtonCotesRule) degree() int {
	if n.open {
		return len(n.coeffs) + 1
	}
	return len(n.coeffs) - 1
}

func (n *NewtonCotesRule) stepSize(a float64, b float64) float64 {
	return (b - a) / float64(n.degree())
}
