package mab

import (
	"context"
	"fmt"

	"github.com/stitchfix/mab/numint"
)

func ExampleBandit_SelectArm() {

	rewards := []Dist{
		Beta(1989, 21290),
		Beta(40, 474),
		Beta(64, 730),
		Beta(71, 818),
		Beta(52, 659),
		Beta(59, 718),
	}

	b := Bandit{
		RewardSource: &RewardStub{Rewards: rewards},
		Strategy:     NewThompson(numint.NewQuadrature()),
		Sampler:      NewSha1Sampler(),
	}

	result, err := b.SelectArm(context.Background(), "12345")
	if err != nil {
		panic(err)
	}
	fmt.Println(result.Arm)

	// Output: 2
}
