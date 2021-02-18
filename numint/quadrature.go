package numint

import (
	"errors"
	"fmt"
	"math"
)

const (
	defaultMaxIter      = 12
	defaultRelTol       = 1e-5
	defaultAbsTol       = 1e-5
	defaultDegree       = 4
	defaultSubIntervals = 2
)

var defaultRule = GaussLegendre(defaultDegree)
var defaultSubdivider = EquallySpaced(defaultSubIntervals)
var defaultTolerance = tolerance{defaultRelTol, defaultAbsTol}

func NewQuadrature(opts ...Option) *Quadrature {
	quad := Quadrature{
		rule:       defaultRule,
		subDivider: defaultSubdivider,
		maxIter:    defaultMaxIter,
		tol:        defaultTolerance,
	}
	for _, opt := range opts {
		opt(&quad)
	}
	return &quad
}

type Rule interface {
	Weights(a float64, b float64) []float64
	Points(a float64, b float64) []float64
}

type SubDivider interface {
	SubDivide([]Interval) []Interval
}

type Integrand func(float64) float64

type Interval struct {
	A float64
	B float64
}

type Quadrature struct {
	rule       Rule
	tol        tolerance
	subDivider SubDivider
	maxIter    int
}

func (q Quadrature) Integrate(f func(float64) float64, a float64, b float64) (float64, error) {
	if !q.canConverge() {
		return math.NaN(), errors.New("integral cannot converge. check tolerance")
	}
	return q.iterativeComposite(f, Interval{a, b})
}

func (q Quadrature) iterativeComposite(f Integrand, interval Interval) (float64, error) {

	intervals := []Interval{interval}

	result, err := q.compositeEstimate(f, intervals)
	if err != nil {
		return math.NaN(), err
	}

	for i := 0; i < q.maxIter; i++ {
		intervals = q.subDivider.SubDivide(intervals)
		prevResult := result
		result, err = q.compositeEstimate(f, intervals)
		if err != nil {
			return math.NaN(), fmt.Errorf(err.Error())
		}
		if q.hasConverged(result, prevResult) {
			return result, nil
		}
	}
	return math.NaN(), errors.New("failed to converge")
}

func (q Quadrature) canConverge() bool {
	return q.tol.absolute > 0 && q.tol.relative > 0
}

func (q Quadrature) hasConverged(result, prevResult float64) bool {
	relErr := relDiff(prevResult, result)
	absErr := absDiff(prevResult, result)

	return relErr <= q.tol.relative && absErr <= q.tol.absolute
}

func (q Quadrature) compositeEstimate(f Integrand, intervals []Interval) (float64, error) {
	total := 0.0
	for i := range intervals {
		result, err := q.singleEstimate(f, intervals[i])
		if err != nil {
			return math.NaN(), err
		}
		total += result
	}
	return total, nil
}

func (q Quadrature) singleEstimate(f Integrand, interval Interval) (float64, error) {

	x := q.rule.Points(interval.A, interval.B)
	w := q.rule.Weights(interval.A, interval.B)

	if len(x) != len(w) {
		return math.NaN(), fmt.Errorf("points and weights must be same length")
	}

	if len(x) == 0 {
		return math.NaN(), fmt.Errorf("points must not be empty")
	}

	sum := 0.0
	for i := range x {
		sum += w[i] * f(x[i])
	}

	return sum, nil
}

type tolerance struct {
	relative float64
	absolute float64
}

func absDiff(a float64, b float64) float64 {
	return math.Abs(a - b)
}

func relDiff(a float64, b float64) float64 {
	switch {
	default:
		return math.Abs(a-b) / b
	case b == 0 && a != 0:
		return math.Abs(a-b) / a
	case b == 0 && a == 0:
		return 0
	}
}
