package mab

import (
	"context"
)

type RewardStub struct {
	Rewards []Dist
}

func (s *RewardStub) GetRewards(context.Context) ([]Dist, error) {
	return s.Rewards, nil
}
