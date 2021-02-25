![Go workflow badge](https://github.com/stitchfix/mab/actions/workflows/go.yml/badge.svg)

# mab: Multi-Armed Bandits Go Library

<p align="center"><img src="https://user-images.githubusercontent.com/5180129/108548622-f2df8200-72a0-11eb-8cc2-b4f1e839dffd.png" width="720"></p>

## Description

### What it is

Mab is a library/framework for scalable and customizable multi-armed bandits. It provides efficient pseudo-random
implementations of epsilon-greedy and Thompson sampling strategies. Arm-selection strategies are decoupled from reward
models, allowing Mab to be used with any reward model whose output can be described as a posterior distribution or point
estimate for each arm.

Mab also provides a numerical one-dimensional integration package, `numint`, which was developed for use by the Mab
Thompson sampler but can also be used as a standalone for numerical integration.

### What it isn't

Mab is not concerned with building, training, or updating bandit reward models. It is focused on efficient pseudo-random
arm selection given the output of a reward model.

## Installation

```
go get -u github.com/stitchfix/mab
```

## Usage

### Bandit

A `Bandit` consists of three components: a `RewardSource`, a `Strategy` and a `Sampler`. Users can provide their own
implementations of each component, or use the Mab implementations.

Example:

```go
package main

import (
	"context"
	"fmt"

	"github.com/stitchfix/mab"
	"github.com/stitchfix/mab/numint"
)

func main() {
	rewards := []mab.Dist{
		mab.Beta(1989, 21290),
		mab.Beta(40, 474),
		mab.Beta(64, 730),
		mab.Beta(71, 818),
		mab.Beta(52, 659),
		mab.Beta(59, 718),
	}

	bandit := mab.Bandit{
		RewardSource: &mab.RewardStub{Rewards: rewards},
		Strategy:     mab.NewThompson(numint.NewQuadrature()),
		Sampler:      mab.NewSha1Sampler(),
	}

	result, err := bandit.SelectArm(context.Background(), "user_id:12345")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Arm)
}
```

`SelectArm` will get the reward estimates from the `RewardSource`, compute arm-selection probabilities using
the `Strategy` and select an arm using the `Sampler`.

The input to `SelectArm` is a `Context`, which can be used to supply request-scoped data to the `RewardSource`
for the purposes of cancellation propagation and/or passing contextual bandit features.

The `unit` input to `SelectArm` is a string that is used for enabling deterministic outcomes. This is useful for
debugging and testing, but can also be used in the context of an experimentation platform to ensure that users get a
consistent experience in between updates to the bandit reward model. Bandits are expected to always provide the same arm
selection for the same set of reward estimates and unit.

The output of `SelectArm` is a struct containing the reward estimates, computed probabilities, and selected arm.

#### RewardSource

A `RewardSource` is expected to provide up-to-date reward estimates for each arm. Users must provide their
own `RewardSource` implementation. Mab only provides a stub implementation for testing and documentation purposes.

```go
type RewardSource interface {
    GetRewards(context.Context) ([]Dist, error)
}
```

A user-defined `RewardSource` is expected to get reward estimates from a database, a cache, or a via HTTP request to a
dedicated reward service. Since a `RewardSource` is likely to require a call to some external service, the `GetRewards`
method includes a `Context` argument. This enables Mab bandits to be used in web services that need to pass
request-scoped context such as request cancellation deadlines.

##### Distributions

Reward estimates are represented as a `Dist` for each arm.

```go
type Dist interface {
    CDF(x float64) float64
    Mean() float64
    Prob(x float64) float64
    Rand() float64
    Support() (float64, float64)
}
```

Mab includes implementations of beta, normal, and point distributions. The beta and normal distributions wrap and
extend [gonum](https://github.com/gonum/gonum/tree/master/stat/distuv) implementations, so they are performant and
reliable.

Mab lets your combine any distribution with any strategy, although some combinations don't make sense in practice. For
example, you could use normal distributions with the epsilon greedy strategy, but the width parameters will just be
ignored. So it might make sense to just use a point distributions for epsilon greedy bandits. Additionally, Mab will let
you use point distributions with a Thompson sampling strategy, but the resulting bandit won't be very useful, since
Thompson sampling relies on distributions having non-zero width.

##### Contextual bandits

For contextual bandits, the `Context` argument can also be used to pass context features to the `RewardSource`.
The `RewardSource` is expected to return the reward estimates conditioned on the context features.

#### Strategy

A Mab `Strategy` computes arm-selection probabilities from the set of reward estimates.

Mab provides the following strategies:

- Thompson sampling
- Epsilon-greedy
- Proportional

##### Thompson sampling

The Thompson sampling strategy computes arm-selection probabilities using the following formula:

![thompson sampling formula](https://user-images.githubusercontent.com/5180129/108559544-4a391e80-72b0-11eb-825c-483aba3dcd18.png)

That is, the probability of selecting an arm under Thompson sampling is the integral of that arm's posterior
PDF times the posterior CDFs of all other arms. The derivation of this formula is left as an exercise for the reader.

Computing these probabilities requires one-dimensional integration, which is provided by the `numint` subpackage.

##### Epsilon-greedy

This is the basic epsilon-greedy selection strategy. The probability of selecting an arm under epsilon greedy is readily
computed from a closed-form solution without the need for numerical integration.

##### Proportional

The proportional sampler computes arm selection probabilities proportional to some input weights. This is not a real
bandit strategy, but exists to allow users to effectively shift the interface between reward sources and bandit
strategies. You can have you `RewardSource` return the desired selection probabilities as `Point` distributions and then
use the `Proportional` strategy to make sure that the sampler uses the provided probabilities directly without an
intermediate probability-computation step.

#### Sampler

A Mab `Sampler` selects an arm given the set of selection probabilities and a string. The default sampler implementation
uses the SHA1 hash of the input string to determine the arm.

### Numerical Integration with numint

The Thompson sampling strategy depends on an integrator for computing probabilities.

```go
type Integrator interface {
    Integrate(f func (float64) float64, a, b float64) (float64, error)
}
```

The `numint` package provides a quadrature-based implementation that can be used for Thompson sampling. It can be used
effectively with just the default settings, or can be fully customized by the user.

The default quadrature rule and other parameters for the `numint` quadrature integrator have been found through
trial-and-error to provide a good tradeoff between speed and reliability for a wide range of inputs including many
combinations of normal and beta distributions.

See the `numint` README and documentation for more details.

## License

Mab and Numint are licensed under the Apache 2.0 license. See the LICENSE file for terms and conditions for use, reproduction, and
distribution.
