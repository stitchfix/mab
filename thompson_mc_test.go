package mab

import (
	"fmt"
)

func ExampleThompsonMC_Probabilities() {
	t := NewThompsonMC(1E6)
	rewards := []Dist{
		Beta(1989, 21290),
		Beta(40, 474),
		Beta(64, 730),
		Beta(71, 818),
		Beta(52, 659),
		Beta(59, 718),
	}
	probs, err := t.ComputeProbs(rewards)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.4f", probs)
	// Output: [0.2967 0.1762 0.2033 0.1687 0.0613 0.0938]
}
