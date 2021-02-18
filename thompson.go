package mab

func NewThompson(integrator Integrator) *Thompson {
	return &Thompson{
		integrator: integrator,
	}
}

type Thompson struct {
	integrator Integrator
	rewards    []Dist
	probs      []float64
}

type Integrator interface {
	Integrate(f func(float64) float64, a, b float64) (float64, error)
}

func (t *Thompson) ComputeProbs(rewards []Dist) ([]float64, error) {
	t.rewards = rewards
	return t.computeProbs()
}

func (t *Thompson) computeProbs() ([]float64, error) {
	t.probs = make([]float64, len(t.rewards))
	for arm := range t.rewards {
		prob, err := t.computeProb(arm)
		if err != nil {
			return nil, err
		}
		t.probs[arm] = prob
	}
	return t.probs, nil
}

func (t *Thompson) computeProb(arm int) (float64, error) {
	integrand := t.integrand(arm)
	xMin, xMax := t.rewards[arm].Support()

	return t.integrator.Integrate(integrand, xMin, xMax)
}

func (t *Thompson) integrand(arm int) func(float64) float64 {
	return func(x float64) float64 {
		total := t.rewards[arm].Prob(x)
		for j := range t.rewards {
			if arm == j {
				continue
			}

			total *= t.rewards[j].CDF(x)
		}
		return total
	}
}
