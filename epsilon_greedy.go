package mab

import (
	"fmt"
	"math"
)

func NewEpsilonGreedy(e float64) *EpsilonGreedy {
	return &EpsilonGreedy{
		Epsilon: e,
	}
}

type EpsilonGreedy struct {
	Epsilon     float64
	meanRewards []float64
}

func (e *EpsilonGreedy) ComputeProbs(rewards []Dist) ([]float64, error) {

	if err := e.validateEpsilon(); err != nil {
		return nil, err
	}

	if len(rewards) == 0 {
		return []float64{}, nil
	}

	e.meanRewards = make([]float64, len(rewards))
	for i, dist := range rewards {
		e.meanRewards[i] = dist.Mean()
	}

	probs := e.computeProbs()
	return probs, nil
}

func (e EpsilonGreedy) computeProbs() []float64 {

	probs := make([]float64, len(e.meanRewards))

	nonNullArms := e.numNonNullArms()
	if nonNullArms == 0 {
		return probs
	}

	maxRewardArmIndices := argsMax(e.meanRewards)
	numMaxima := len(maxRewardArmIndices)

	for i := range e.meanRewards {
		if isIn(maxRewardArmIndices, i) {
			probs[i] = (1-e.Epsilon)/float64(numMaxima) + e.Epsilon/float64(nonNullArms)
		} else {
			if math.IsInf(e.meanRewards[i], -1) {
				probs[i] = 0
			} else {
				probs[i] = e.Epsilon / float64(nonNullArms)
			}
		}
	}

	return probs
}

func (e EpsilonGreedy) numNonNullArms() int {
	count := 0
	for _, val := range e.meanRewards {
		if val > math.Inf(-1) {
			count += 1
		}
	}
	return count
}

func (e EpsilonGreedy) validateEpsilon() error {
	if e.Epsilon < 0 || e.Epsilon > 1 {
		return fmt.Errorf("invalid Epsilon value: %v. Must be between 0 and 1", e.Epsilon)
	}
	return nil
}

func isIn(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func argsMax(vals []float64) []int {
	var maxVal = math.Inf(-1)
	var maxArgs []int
	for i, val := range vals {
		if val > maxVal {
			maxArgs = []int{i}
			maxVal = val
		} else if val == maxVal {
			maxArgs = append(maxArgs, i)
		}
	}
	return maxArgs
}
