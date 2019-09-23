package main

import (
	"math/rand"
)

const (
	next_prime int = 2147483587
	max_value int = next_prime - 1
)

type MinHasher struct {
	coeffA []int
	coeffB []int
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
func (mh* MinHasher) Hash(doc map[int]bool) (sigs []int) {
	sigs = make([]int, mh.num_hashes)
	for i := 0; i < mh.num_hashes; i++ {
		min := next_prime + 1
		for shingle := range doc {
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

func pickRandCoeffs(k int) (result []int) {
	result = make([]int, k)
	var seen map[int]bool

	seen = make(map[int]bool, k)
	i := 0
	for k > 0 {
		randIndex := rand.Intn(max_value)
		for seen[randIndex] {
			randIndex = rand.Intn(max_value)
		}
		result[i] = randIndex
		seen[randIndex] = true
		k--
		i++
	}

	return result
}

	
