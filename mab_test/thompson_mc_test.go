package mab

import (
	"fmt"

	"github.com/stitchfix/mab"
)

func ExampleThompsonMC_ComputeProbs() {
	t := mab.NewThompsonMC(1E6)
	rewards := []mab.Dist{
		mab.Beta(1989, 21290),
		mab.Beta(40, 474),
		mab.Beta(64, 730),
		mab.Beta(71, 818),
		mab.Beta(52, 659),
		mab.Beta(59, 718),
	}
	probs, err := t.ComputeProbs(rewards)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.4f", probs)
	// Output: [0.2967 0.1762 0.2033 0.1687 0.0613 0.0938]
}
