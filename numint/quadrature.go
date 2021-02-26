// Package numint provides rules and methods for one-dimensional numerical quadrature
package numint

import (
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

// NewQuadrature returns a pointer to new Quadrature with any Option arguments applied.
// For example:
// 	q := NewQuadrature()
// Returns a Quadrature with all default settings.
// The default settings can be overridden with Option functions:
//	q := NewQuadrature(WithRule(GaussLegendre(4), WithMaxIter(10), WithRelTol(0.01))
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

// Rule is an interface that provides Weights and sampling Points to be used during numerical quadrature.
type Rule interface {
	Weights(a float64, b float64) []float64
	Points(a float64, b float64) []float64
}

// SubDivider determines how to sub-divide intervals for each iteration of numerical quadrature.
// It takes a slice of intervals and returns the a flat slice containing all the sub-divided intervals.
// For example, using the EquallySpaced subdivider to divide the intervals [0, 1], [1, 2] into two equally-spaced intervals each:
// 	sub := EquallySpaced(2)
// 	result := sub.SubDivide([]Interval{{0, 1}, {1, 2}})
// Results in the sub-intervals:
// 	[]Interval{{0, 0.5}, {0.5, 1}, {1, 1.5}, {1.5, 2}}
type SubDivider interface {
	SubDivide([]Interval) []Interval
}

type integrand func(float64) float64

// Interval represents a finite interval between A and B, where B > A.
type Interval struct {
	A float64
	B float64
}

// Quadrature contains the rule, subdivider, tolerance, and max iterations for numerical quadrature.
// These fields can all be specified using NewQuadrature with the corresponding option functions.
type Quadrature struct {
	rule       Rule
	tol        tolerance
	subDivider SubDivider
	maxIter    int
}

// Integrate computes an estimate of the integral of f from a to b.
// It works by first getting the Points and Weights from the Rule for the interval [a, b]
// then computing the sum of w_i * f(x_i) where w_i are the weights and p_i are the points.
// The next step is to subdivide the original interval [a, b] using the SubDivider,
// then compute the same estimate summed over the sub-intervals.
// This process is repeated until the absolute difference between successive iterations is less than the specified tolerance,
// or until maxIter is reached.
// If absolute tolerance is set using WithAbsTol, only absolute tolerance is checked.
// If relative tolerance is set using WithRelTol, only relative tolerance is checked.
// If both absolute and relative tolerances are set using WithAbsAndRelTol, then the absolute difference must be less than *both* tolerances for the algorithm to converge.
// If the max iteration threshold is reached without reaching the specified tolerance, Integrate returns the final result and an error.
// The max iteration threshold can be specified using WithMaxIter as an argument to NewQuadrature.
func (q Quadrature) Integrate(f func(float64) float64, a float64, b float64) (float64, error) {
	if a == b {
		return 0, nil
	}
	if !q.canConverge() {
		return math.NaN(), fmt.Errorf("integral cannot converge. check tolerance")
	}
	return q.iterativeComposite(f, Interval{a, b})
}

func (q Quadrature) iterativeComposite(f integrand, interval Interval) (float64, error) {

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
			return result, err
		}
		if q.hasConverged(result, prevResult) {
			return result, nil
		}
	}
	return result, fmt.Errorf("failed to converge")
}

func (q Quadrature) canConverge() bool {
	return q.tol.absolute > 0 && q.tol.relative > 0
}

func (q Quadrature) hasConverged(result, prevResult float64) bool {
	relErr := relDiff(prevResult, result)
	absErr := absDiff(prevResult, result)

	return relErr <= q.tol.relative && absErr <= q.tol.absolute
}

func (q Quadrature) compositeEstimate(f integrand, intervals []Interval) (float64, error) {
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

func (q Quadrature) singleEstimate(f integrand, interval Interval) (float64, error) {

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
