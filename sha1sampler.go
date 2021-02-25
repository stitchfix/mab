package mab

import (
	"crypto/sha1"
	"fmt"
	"strconv"
)

const defaultNumBuckets = 1000

func NewSha1Sampler() *Sha1Sampler {
	return &Sha1Sampler{
		numBuckets: defaultNumBuckets,
	}
}

// Sha1Sampler is a Sampler that uses the SHA1 hash of input unit to select an arm index with probability proportional to some given weights.
type Sha1Sampler struct {
	numBuckets int
}

// Sample returns the selected arm for a given set of weights and input unit.
// An error is returned if any negative weight is encountered.
func (s *Sha1Sampler) Sample(weights []float64, unit string) (int, error) {

	checkSum := sha1.Sum([]byte(unit))

	hexDigest := fmt.Sprintf("%x", checkSum[0:8])
	hexDigest = hexDigest[0 : len(hexDigest)-1]

	uBucket64, _ := strconv.ParseUint(hexDigest, 16, 64)

	bucket := int(uBucket64) % s.numBuckets

	return s.getIndex(weights, bucket)
}

func (s *Sha1Sampler) sum(weights []float64) float64 {
	sum := 0.0
	for _, w := range weights {
		sum += w
	}
	return sum
}

func (s *Sha1Sampler) getIndex(weights []float64, bucket int) (int, error) {
	sumWeights := s.sum(weights)
	if sumWeights <= 0 {
		return -1, fmt.Errorf("sum(weights) must be positive. got=%0.2f", sumWeights)
	}

	curBucket := -1.0

	for i, w := range weights {
		if w < 0 {
			return -1, fmt.Errorf("negative weight")
		}
		curBucket += w * float64(s.numBuckets) / sumWeights
		if curBucket >= float64(bucket) {
			return i, nil
		}
	}

	return -1, fmt.Errorf("bucket out of range") // this code should be unreachable
}
