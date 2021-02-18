package mab

func NewThompsonMC(numIterations int) *ThompsonMC {
	return &ThompsonMC{
		NumIterations: numIterations,
	}
}

type ThompsonMC struct {
	NumIterations   int
	rewards         []Dist
	counts, samples []float64
}

func (t *ThompsonMC) ComputeProbs(rewards []Dist) ([]float64, error) {
	t.rewards = rewards
	return t.computeProbs(), nil
}

func (t *ThompsonMC) computeProbs() []float64 {
	t.counts = make([]float64, len(t.rewards))
	t.samples = make([]float64, len(t.rewards))

	for i := 0; i < t.NumIterations; i++ {
		for j, r := range t.rewards {
			t.samples[j] = r.Rand()
		}

		maxArgs := argsMax(t.samples)
		numMaxima := len(maxArgs)

		for _, maximum := range maxArgs {
			t.counts[maximum] += 1.0 / float64(numMaxima)
		}
	}

	result := make([]float64, len(t.rewards))

	for i := range t.rewards {
		result[i] = t.counts[i] / float64(t.NumIterations)
	}

	return result
}
