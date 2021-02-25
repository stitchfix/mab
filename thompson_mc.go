package mab

// NewThompsonMC returns a new ThompsonMC with numIterations.
func NewThompsonMC(numIterations int) *ThompsonMC {
	return &ThompsonMC{
		NumIterations: numIterations,
	}
}

// ThompsonMC is a Monte-Carlo based implementation of Thompson sampling Strategy.
// It should not be used in production but is provided only as an example and for comparison with the Thompson Strategy,
// which is much faster and more accurate.
type ThompsonMC struct {
	NumIterations   int
	rewards         []Dist
	counts, samples []float64
}

// ComputeProbs estimates the arm-selection probabilities by repeatedly sampling from the Dist for each arm,
// and counting how many times each arm yields the maximal sampled value.
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
