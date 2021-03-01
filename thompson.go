package mab

import (
	"sync"
)

func NewThompson(integrator Integrator) *Thompson {
	return &Thompson{
		integrator: integrator,
	}
}

type Thompson struct {
	integrator Integrator
}

type Integrator interface {
	Integrate(f func(float64) float64, a, b float64) (float64, error)
}

func (t *Thompson) ComputeProbs(rewards []Dist) ([]float64, error) {
	if len(rewards) == 0 {
		return []float64{}, nil
	}

	integrals := t.integrals(rewards)
	return t.integrateParallel(integrals)
}

type integral struct {
	integrand integrand
	interval  interval
}

type integrand func(float64) float64
type interval struct{ a, b float64 }

func (t *Thompson) integrals(rewards []Dist) []integral {
	result := make([]integral, len(rewards))
	for i := range rewards {
		result[i].integrand = t.integrand(rewards, i)
		result[i].interval.a, result[i].interval.b = rewards[i].Support()
	}
	return result
}

func (t *Thompson) integrand(rewards []Dist, arm int) integrand {
	return func(x float64) float64 {
		total := rewards[arm].Prob(x)
		for j := range rewards {
			if arm == j {
				continue
			}

			total *= rewards[j].CDF(x)
		}
		return total
	}
}

func (t *Thompson) integrateParallel(integrals []integral) ([]float64, error) {
	n := len(integrals)

	results := make([]float64, n)
	errs := make([]error, n)

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int, xi integral) {
			results[i], errs[i] = t.integrator.Integrate(xi.integrand, xi.interval.a, xi.interval.b)
			wg.Done()
		}(i, integrals[i])
	}

	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}
