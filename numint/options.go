package numint

import "math"

type Option func(*Quadrature)

func WithMaxIter(m int) Option {
	return func(q *Quadrature) {
		q.maxIter = m
	}
}

func WithAbsTol(absTol float64) Option {
	return func(q *Quadrature) {
		q.tol = tolerance{
			absolute: absTol,
			relative: math.Inf(1),
		}
	}
}

func WithRelTol(relTol float64) Option {
	return func(q *Quadrature) {
		q.tol = tolerance{
			absolute: math.Inf(1),
			relative: relTol,
		}
	}
}

func WithAbsAndRelTol(absTol float64, relTol float64) Option {
	return func(q *Quadrature) {
		q.tol = tolerance{
			absolute: absTol,
			relative: relTol,
		}
	}
}

func WithRule(rule Rule) Option {
	return func(q *Quadrature) {
		q.rule = rule
	}
}

func WithSubDivider(s SubDivider) Option {
	return func(q *Quadrature) {
		q.subDivider = s
	}
}
