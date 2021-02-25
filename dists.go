package mab

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/stat/distuv"
)

// Normal is a normal distribution for use with any bandit strategy.
// For the purposes of Thompson sampling, it is truncated at mean +/- 4*sigma
func Normal(mu, sigma float64) NormalDist {
	return NormalDist{distuv.Normal{Mu: mu, Sigma: sigma}}
}

type NormalDist struct {
	distuv.Normal
}

func (n NormalDist) Support() (float64, float64) {
	width := 4.0
	return n.Mu - width*n.Sigma, n.Mu + width*n.Sigma
}

func (n NormalDist) String() string {
	return fmt.Sprintf("Normal(%f,%f)", n.Mu, n.Sigma)
}

// Beta is a beta distribution for use with any bandit strategy.
func Beta(alpha, beta float64) BetaDist {
	return BetaDist{distuv.Beta{Alpha: alpha, Beta: beta}}
}

type BetaDist struct {
	distuv.Beta
}

func (b BetaDist) Support() (float64, float64) {
	return 0, 1
}

func (b BetaDist) String() string {
	return fmt.Sprintf("Beta(%f,%f)", b.Beta.Alpha, b.Beta.Beta)
}

// Point is used for reward models that just provide point estimates. Don't use with Thompson sampling.
func Point(mu float64) PointDist {
	return PointDist{mu}
}

type PointDist struct {
	Mu float64
}

func (p PointDist) Mean() float64 {
	return p.Mu
}

func (p PointDist) CDF(x float64) float64 {
	if x >= p.Mu {
		return 1
	}
	return 0
}

func (p PointDist) Prob(x float64) float64 {
	if x == p.Mu {
		return math.NaN()
	}
	return 0
}

func (p PointDist) Rand() float64 {
	return p.Mu
}

func (p PointDist) Support() (float64, float64) {
	return p.Mu, p.Mu
}

func (p PointDist) String() string {
	if math.IsInf(p.Mu, -1) {
		return "Null()"
	}
	return fmt.Sprintf("Point(%f)", p.Mu)
}

// Null returns a PointDist with mean equal to negative infinity. This is a special value that indicates
// to a Strategy that this arm should get selection probability zero.
func Null() PointDist {
	return PointDist{math.Inf(-1)}
}
