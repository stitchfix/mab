package mab_test

import (
	"context"
	"testing"

	"github.com/stitchfix/mab"
	"github.com/stitchfix/mab/numint"
)

func TestThompson_SelectArm(t *testing.T) {

	rewards := []mab.Dist{
		mab.Beta(1989, 21290),
		mab.Beta(40, 474),
		mab.Beta(64, 730),
		mab.Beta(71, 818),
		mab.Beta(52, 659),
		mab.Beta(59, 718),
	}

	b := mab.Bandit{
		RewardSource: &mab.RewardStub{Rewards: rewards},
		Strategy:     mab.NewThompson(numint.NewQuadrature()),
		Sampler:      mab.NewSha1Sampler(),
	}

	result, err := b.SelectArm(context.Background(), "12345", nil)
	if err != nil {
		t.Fatal(err)
	}

	actual := result.Arm
	expected := 2

	if actual != expected {
		t.Errorf("result not %d, got=%d", expected, actual)
	}
}

func TestEpsilon_SelectArm(t *testing.T) {
	rewards := map[string][]mab.Dist{
		"blue": {mab.Point(-0.5), mab.Point(0.5)},
		"red":  {mab.Point(0.5), mab.Point(-0.5)},
	}

	b := mab.Bandit{
		RewardSource: &mab.ContextualRewardStub{Rewards: rewards},
		Strategy:     &mab.EpsilonGreedy{Epsilon: 0.0},
		Sampler:      mab.NewSha1Sampler(),
	}

	result, err := b.SelectArm(context.Background(), "12345", "red")
	if err != nil {
		t.Fatal(err)
	}

	actual := result.Arm
	expected := 0

	if actual != expected {
		t.Errorf("result not %d, got=%d", expected, actual)
	}
}
