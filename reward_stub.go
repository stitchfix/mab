package mab

import (
	"context"
	"fmt"
)

// RewardStub is a static non-contextual RewardSource that can be used for testing and development.
type RewardStub struct {
	Rewards []Dist
}

// GetRewards gets the static rewards
func (s *RewardStub) GetRewards(context.Context, interface{}) ([]Dist, error) {
	return s.Rewards, nil
}

// ContextualRewardStub is a static contextual RewardSource that can be used for testing and development of contextual bandits.
// It assumes that the context can be specified with a string.
type ContextualRewardStub struct {
	Rewards map[string][]Dist
}

// GetRewards gets the static rewards for a given banditContext string.
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
