package lib

import (
	"math/rand"
)

const (
	next_prime uint32 = 2147483587
	max_value uint32 = next_prime - 1
)

type MinHasher struct {
	coeffA []uint32
	coeffB []uint32
	num_hashes int
}

func NewMinhash(num_hashes int) *MinHasher {
	mh := new(MinHasher)
	mh.num_hashes = num_hashes
	mh.coeffA = pickRandCoeffs(num_hashes)
	mh.coeffB = pickRandCoeffs(num_hashes)
	return mh
}

// The input here is a map of integers to bool.  Each integer corresponds to a feature
// ID.  Users of this code will typically take shingles and reduce them with the hashing
// trick.
func (mh* MinHasher) Hash(doc []uint32) (sigs []uint32) {
	sigs = make([]uint32, mh.num_hashes)
	for i := 0; i < mh.num_hashes; i++ {
		min := next_prime + 1
		for _, shingle := range doc {
			shingle = shingle % max_value
			h := (mh.coeffA[i] * shingle + mh.coeffB[i]) % next_prime
			if h < min {
				min = h
			}
		}
		sigs[i] = min;
	}
	return sigs
}

func pickRandCoeffs(k int) (result []uint32) {
	result = make([]uint32, k)
	var seen map[uint32]bool

	seen = make(map[uint32]bool, k)
	i := 0
	for k > 0 {
		randIndex := rand.Uint32() % max_value
		for seen[randIndex] {
			randIndex = rand.Uint32() % max_value
		}
		result[i] = randIndex
		seen[randIndex] = true
		k--
		i++
	}

	return result
}

	
