# Mab
Multi-Armed Bandits Go Library

<p align="center"><img src="https://user-images.githubusercontent.com/5180129/108548622-f2df8200-72a0-11eb-8cc2-b4f1e839dffd.png" width="720"></p>
<p align="center">
	<a href="https://github.com/stitchfix/mab/actions/workflows/go.yml"><img src="https://github.com/stitchfix/mab/actions/workflows/go.yml/badge.svg" alt="Build Status"></img></a>
	<a href="https://goreportcard.com/report/github.com/stitchfix/mab"><img src="https://goreportcard.com/badge/github.com/stitchfix/mab" alt="Go Report Card"></img></a>
	<a href="https://pkg.go.dev/github.com/stitchfix/mab"><img src="https://pkg.go.dev/badge/github.com/stitchfix/mab.svg" alt="Go Reference"></img></a>
</p>

* [Description](#description)
* [Installation](#installation)
* [Usage](#usage)
  + [Creating a bandit and selecting arms](#bandit)
  + [Numerical integration with `numint`](#numint)
* [Documentation](#documentation)
* [License](#license)

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

	rewards := map[string][]mab.Dist{
		"us": {
			mab.Beta(40, 474),
			mab.Beta(64, 730),
			mab.Beta(71, 818),
		},
		"uk": {
			mab.Beta(25, 254),
			mab.Beta(100, 430),
			mab.Beta(30, 503),
		},
	}

	bandit := mab.Bandit{
		RewardSource: &mab.ContextualRewardStub{rewards},
		Strategy:     mab.NewThompson(numint.NewQuadrature()),
		Sampler:      mab.NewSha1Sampler(),
	}

	result, err := bandit.SelectArm(context.Background(), "user_id:12345", "us")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}
```

`SelectArm` will get the reward estimates from the `RewardSource`, compute arm-selection probabilities using
the `Strategy` and select an arm using the `Sampler`.

There is an unfortunate name collision between Go's `context.Context` type and the context a contextual bandit.
In Mab, the `context.Context` variables will always be named `ctx`, while the variables used for bandit context will be called `banditContext`.

Go's `context.Context` should be used to pass request-scoped data to the RewardSource, and it is best practice to only use it for cancellation propagation or passing non-controlling data such as request IDs.

The values needed by the contextual bandit to determine the reward estimates should be passed using the last argument, which is named `banditContext`.

The `unit` input to `SelectArm` is a string that is used for enabling deterministic outcomes. This is useful for
debugging and testing, but can also be used to ensure that users get a consistent experience in between updates to the bandit reward model.
Bandits are expected to always provide the same arm selection for the same set of reward estimates and input unit string.

The output of `SelectArm` is a struct containing the reward estimates, computed probabilities, and selected arm.

#### RewardSource

A `RewardSource` is expected to provide up-to-date reward estimates for each arm, given some context data.
Mab provides a basic implementation (`HTTPSource`) that can be used for requesting rewards from an HTTP service, and some stubs that can be used for testing and development.

```go
type RewardSource interface {
    GetRewards(context.Context, interface{}) ([]Dist, error)
}
```

A typical `RewardSource` implementation is expected to get reward estimates from a database, a cache, or a via HTTP request to a
dedicated reward service. Since a `RewardSource` is likely to require a call to some external service, the `GetRewards`
method includes a `context.Context`-type argument. This enables Mab bandits to be used in web services that need to pass
request-scoped data such as request timeouts and cancellation propagation. The second argument should be used to pass bandit context data to the reward source.
The reward source must return one distribution per arm, conditional on the bandit context.

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

Mab lets your combine any distribution with any strategy, although some combinations don't make sense in practice. 

For epsilon greedy, you will most likely use `Point` distributions, since the algorithm only cares about the mean of the reward estimate.
Other distributions can be used, as long as they implement a `Mean()` that returns well-defined values.

For Thompson sampling, it is recommended to use `Normal` or `Beta` distributions. Since Thompson sampling is based on sampling from finite-width distributions, you won't get a useful bandit by using `Point` distributions with the `Thompson` strategy.

The `Null()` function returns a `Point` distribution at negative infinity (`math.Inf(-1)`). This indicates to the `Strategy` that this arm should never be selected. Each `Strategy` must account for any number of Null distributions and return zero probability for the null arms and the correct set of probabilities for the non-null arms, as if the null arms were not present.

#### Strategy

A Mab `Strategy` computes arm-selection probabilities from the set of reward estimates.

Mab provides the following strategies:

- Thompson sampling (`mab.Thompson`)
- Epsilon-greedy (`mab.EpsilonGreedy`)
- Proportional (`mab.Proportional`)

Mab also provides a Monte-Carlo based Thompson-sampling strategy (`mab.ThompsonMC`) but it is much slower an less accurate than `mab.Thompson`, which is based on numerical integration. It is not recommended to use `ThompsonMC` in production.

##### Thompson sampling

The Thompson sampling strategy computes arm-selection probabilities using the following formula:

![thompson sampling formula](https://user-images.githubusercontent.com/5180129/108559544-4a391e80-72b0-11eb-825c-483aba3dcd18.png)

That is, the probability of selecting an arm under Thompson sampling is the integral of that arm's posterior
PDF times the posterior CDFs of all other arms. The derivation of this formula is left as an exercise for the reader.

Computing these probabilities requires one-dimensional integration, which is provided by the `numint` subpackage.

The limits of integration are determined by the `Support` of the arms' distribution, so `Point` distributions will always get zero probability using Thompson sampling.

##### Epsilon-greedy

This is the basic epsilon-greedy selection strategy. The probability of selecting an arm under epsilon greedy is readily
computed from a closed-form solution without the need for numerical integration. It is based on the `Mean` of the reward estimate.

##### Proportional

The proportional sampler computes arm selection probabilities proportional to some input weights. This is not a real
bandit strategy, but exists to allow users to effectively shift the interface between reward sources and bandit
strategies. You can create a `RewardSource` that returns the desired selection weights as `Point` distributions and then
use the `Proportional` strategy to make sure that the sampler uses the normalized weights as the probability distribution for arm selection.

#### Sampler

A Mab `Sampler` selects an arm given the set of selection probabilities and a string. The default sampler implementation
uses the SHA1 hash of the input string (mod 1000) to determine the arm.

### Numint

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

## Documentation

More detailed refence docs can be found on [pkg.go.dev](https://pkg.go.dev/github.com/stitchfix/mab)

## License

Mab and Numint are licensed under the Apache 2.0 license. See the LICENSE file for terms and conditions for use, reproduction, and
distribution.
