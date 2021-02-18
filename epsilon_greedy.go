package mab

import (
	"fmt"
	"math"
)

type EpsilonGreedy struct {
	Epsilon     float64
	meanRewards []float64
}

func (e *EpsilonGreedy) ComputeProbs(rewards []Dist) ([]float64, error) {

	if err := e.validateEpsilon(); err != nil {
		return nil, err
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

	maxRewardArmIndices := argsMax(e.meanRewards)
	numMaxima := len(maxRewardArmIndices)
	numArms := len(e.meanRewards)

	for i := range e.meanRewards {
		if isIn(maxRewardArmIndices, i) {
			probs[i] = (1-e.Epsilon)/float64(numMaxima) + e.Epsilon/float64(numArms)
		} else {
			probs[i] = e.Epsilon / float64(len(e.meanRewards))
		}
	}

	return probs
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
