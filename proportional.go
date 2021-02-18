package mab

import "fmt"

type Proportional struct {
	meanRewards, probs []float64
}

func (p *Proportional) ComputeProbs(rewards []Dist) ([]float64, error) {

	p.meanRewards = make([]float64, len(rewards))
	for i, dist := range rewards {
		mean := dist.Mean()
		if mean < 0 {
			return nil, fmt.Errorf("negative mean reward")
		}

		p.meanRewards[i] = mean
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
	for i, mean := range p.meanRewards {
		p.probs[i] = mean / norm
	}

	return p.probs, nil
}
