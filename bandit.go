package mab

import (
	"context"
)

// A Bandit gets reward values from a RewardSource, computes selection probabilities using a Strategy, and selects
// an arm using a Sampler.
type Bandit struct {
	RewardSource
	Strategy
	Sampler
}

// SelectArm gets the current reward estimates, computes the arm selection probabilities, and selects and arm index.
// Returns a partial result and an error message if an error is encountered at any point.
// For example, if the reward estimates were retrieved, but an error was encountered during the probability computation,
// the result will contain the reward estimates, but no probabilities or arm index.
// There is an unfortunate name collision between a multi-armed bandit context and Go's context.Context type.
// The context.Context argument should only be used for passing request-scoped data to an external reward service, such
// as timeouts and cancellation propagation.
// The banditContext argument is used to pass bandit context features to the reward source for contextual bandits.
// The unit argument is a string that will be hashed to select an arm with the pseudo-random sampler.
// SelectArm is deterministic for a fixed unit and set of reward estimates from the RewardSource.
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

// Result is the return type for a call to Bandit.SelectArm.
// It will contain the reward estimates provided by the RewardSource, the computed arm selection probabilities,
// and the index of the selected arm.
type Result struct {
	Rewards []Dist
	Probs   []float64
	Arm     int
}

// A Dist represents a one-dimensional probability distribution.
// Reward estimates are represented as a Dist for each arm.
// Strategies compute arm-selection probabilities using the Dist interface.
// This allows for combining different distributions with different strategies.
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

// A RewardSource provides the current reward estimates, in the form of a Dist for each arm.
// There is an unfortunate name collision between a multi-armed bandit context and Go's Context type.
// The first argument is a context.Context and should only be used for passing request-scoped data to an external reward service.
// If the RewardSource does not require an external request, this first argument should always be context.Background()
// The second argument is used to pass context values to the reward source for contextual bandits.
// A RewardSource implementation should provide the reward estimates conditioned on the value of banditContext.
// For non-contextual bandits, banditContext can be nil.
type RewardSource interface {
	GetRewards(ctx context.Context, banditContext interface{}) ([]Dist, error)
}

// A Strategy computes arm selection probabilities from a slice of Distributions.
type Strategy interface {
	ComputeProbs([]Dist) ([]float64, error)
}

// A Sampler returns a pseudo-random arm index given a set of probabilities and a string to hash.
// Samplers should always return the same arm index for the same set of probabilities and unit value.
type Sampler interface {
	Sample(probs []float64, unit string) (int, error)
}
