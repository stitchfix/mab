package mab

import (
	"fmt"

	"github.com/stitchfix/mab/numint"
)

func ExampleThompson_Probabilities() {
	t := NewThompson(numint.NewQuadrature())
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
	// Output: [0.2963 0.1760 0.2034 0.1690 0.0614 0.0939]
}
