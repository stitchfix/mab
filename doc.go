// Copyright Stitch Fix, Inc. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package mab provides pseudo-random multi-armed bandit selection strategies.
A Bandit consists of a reward source that provides the current reward estimate for each arm, a strategy for computing
arm-selection probabilities, and a sampler for pseudo-random arm selection from a set of probabilities.

Mab is not concerned with building, training, or updating reward models. It is only used for selecting arms
given the output of some reward model.

Distributions

Reward estimates are provided as a set of distributions, one for each arm. Any type that satisfies the Dist interface
can be used as a reward estimate for an arm. The mab package has three implementations: Normal, Beta, and Point.
Any distribution can be used with any Strategy, although some combinations don't make sense in practice.
For example, Thompson sampling requires finite-width distributions, so it wouldn't make sense to provide Point distributions
to a Thompson sampling policy. Mab does not require that all arms use the same distribution type.

Reward Sources

A RewardSource provides current reward estimates. Users are expected to implement their own RewardSources.
Mab only provides a RewardStub implementation for testing.
A production implementation of RewardSource will likely get current reward estimates from a cache, a database,
or via a web request to a reward service.

Strategies

A Strategy computes arm selection probabilities from a set of reward estimates. Example Strategies are Thompson sampling
and EpsilonGreedy. The ThompsonMC strategy is provided only for comparison to the much faster integration-based Thompson strategy,
and is far too slow to use in an online system. The Proportional strategy results in arm-selection probability proportional
to the mean reward for each arm.

Samplers

A Sampler selects an arm given a set of selection probabilities. Mab provides a SHA1-based sampler, which selects
an arm based on hashing an input string. This ensures that repeated calls with the same set of probabilities and same unit
result in the same arm selection.

Arms and Parameters

An arm is a set of parameters that can be selected. A parameter has a name, which must be a string, and a value, which
can be of any type. So each arm is represented as a map[string]interface{}.
*/
package mab
