package mab

import (
	"context"
	"fmt"
)

type RewardStub struct {
	Rewards []Dist
}

func (s *RewardStub) GetRewards(context.Context, interface{}) ([]Dist, error) {
	return s.Rewards, nil
}

type ContextualRewardStub struct {
	Rewards map[string][]Dist
}

func (c *ContextualRewardStub) GetRewards(ctx context.Context, banditContext interface{}) ([]Dist, error) {
	key, ok := banditContext.(string)
	if !ok {
		return nil, fmt.Errorf("banditContext must be a string")
	}

	val, ok := c.Rewards[key]

	if !ok {
		return nil, fmt.Errorf("no distributions for %s", val)
	}

	return val, nil
}
