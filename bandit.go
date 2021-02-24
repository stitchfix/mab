package mab

import (
	"context"
)

// Bandit gets reward values from a RewardSource, computes selection probabilities using a Strategy, and selects
// an arm using a Sampler.
type Bandit struct {
	RewardSource
	Strategy
	Sampler
}

type Result struct {
	Rewards []Dist
	Probs   []float64
	Arm     int
}

// SelectArm gets the current reward estimates, computes the arm selection probability, and selects and arm index.
func (b *Bandit) SelectArm(ctx context.Context, unit string, banditContext interface{}) (Result, error) {

	res := Result{
		Rewards: make([]Dist, 0),
		Probs:   make([]float64, 0),
		Arm:     -1,
	}

	rewards, err := b.GetRewards(ctx, banditContext)
	if err != nil {
		return res, err
	}

	res.Rewards = rewards

	probs, err := b.ComputeProbs(rewards)
	if err != nil {
		return res, err
	}

	res.Probs = probs

	result, err := b.Sample(probs, unit)
	if err != nil {
		return res, err
	}

	res.Arm = result

	return res, nil
}

// RewardSource provides the current reward estimates, as a Dist for each arm.
// Features can be passed to the RewardSource using the Context argument, which is useful for contextual bandits.
// The RewardSource should provide the reward estimates conditioned on those context features.
type RewardSource interface {
	GetRewards(context.Context, interface{}) ([]Dist, error)
}

// Dist represents a one-dimensional probability distribution.
type Dist interface {
	// CDF returns the cumulative distribution function evaluated at x.
	CDF(x float64) float64

	// Mean returns the mean of the distribution.
	Mean() float64

	// Prob returns the probability density function or probability mass function evaluated at x.
	Prob(x float64) float64

	// Rand returns a pseudo-random sample drawn from the distribution.
	Rand() float64

	// Support returns the range of values over which the distribution is considered non-zero for the purposes of numerical integration.
	Support() (float64, float64)
}

// Strategy computes arm selection probabilities from a slice of Distributions.
// The output probabilities slice should be the same length as the input Dist slice.
type Strategy interface {
	ComputeProbs([]Dist) ([]float64, error)
}

// Sampler returns a pseudo-random arm index given a set of probabilities and a unit.
// Samplers should always return the same arm index for the same set of probabilities and unit.
type Sampler interface {
	Sample(probs []float64, unit string) (int, error)
}
