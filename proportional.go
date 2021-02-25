package mab

import (
	"fmt"
	"math"
)

func NewProportional() *Proportional {
	return &Proportional{}
}

type Proportional struct {
	meanRewards, probs []float64
}

func (p *Proportional) ComputeProbs(rewards []Dist) ([]float64, error) {

	p.meanRewards = make([]float64, len(rewards))
	for i, dist := range rewards {
		mean := dist.Mean()

		switch {
		default:
			p.meanRewards[i] = mean
		case mean > math.Inf(-1) && mean < 0:
			return nil, fmt.Errorf("negative mean reward")
		case math.IsInf(mean, -1): // indicates a NULL distribution
			p.meanRewards[i] = 0
		}
	}

	return p.computeProbs()
}

func (p Proportional) computeProbs() ([]float64, error) {
	norm := 0.0
	for _, r := range p.meanRewards {
		if r < 0 {
			return nil, fmt.Errorf("negative mean reward: %+v", r)
		}
		norm += r
	}

	p.probs = make([]float64, len(p.meanRewards))

	if norm == 0 {
		return p.probs, nil
	}

	for i, mean := range p.meanRewards {
		p.probs[i] = mean / norm
	}

	return p.probs, nil
}
