package numint

import "math"

// Option is a function that can be passed to NewQuadrature to override the default settings.
type Option func(*Quadrature)

// WithMaxIter sets the max iterations to m
func WithMaxIter(m int) Option {
	return func(q *Quadrature) {
		q.maxIter = m
	}
}

// WithAbsTol sets the absolute tolerance convergence criteria to absTol and sets the relative tolerance to be ignored.
func WithAbsTol(absTol float64) Option {
	return func(q *Quadrature) {
		q.tol = tolerance{
			absolute: absTol,
			relative: math.Inf(1),
		}
	}
}

// WithRelTol sets the relative tolerance convergence criteria to relTol and sets the absolute tolerance to be ignored.
func WithRelTol(relTol float64) Option {
	return func(q *Quadrature) {
		q.tol = tolerance{
			absolute: math.Inf(1),
			relative: relTol,
		}
	}
}

// WithAbsAndRelTol sets both the absolute and relative tolerances so that the absolute difference and relative differences
// between successive iterations must both meet a threshold for convergence.
func WithAbsAndRelTol(absTol float64, relTol float64) Option {
	return func(q *Quadrature) {
		q.tol = tolerance{
			absolute: absTol,
			relative: relTol,
		}
	}
}

// WithRule sets the rule that should be used for each iteration of numerical quadrature.
func WithRule(rule Rule) Option {
	return func(q *Quadrature) {
		q.rule = rule
	}
}

// WithSubDivider sets the subdivider that should be used to compute the set of sub-intervals for each iteration.
func WithSubDivider(s SubDivider) Option {
	return func(q *Quadrature) {
		q.subDivider = s
	}
}
