package main

import (
	"testing"
	"math/rand"
)

func TestSetup(t *testing.T) {
	mh := NewMinhash(256)
	for i, c := range mh.coeffA {
		if c < 0 {
			t.Errorf("coeffA[%d] is negative: %q", i, c)
		}
	}	
	for i, c := range mh.coeffB {
		if c < 0 {
			t.Errorf("coeffB[%d] is negative: %q", i, c)
		}
	}
}

func TestHash(t *testing.T) {
	num_terms := 1234
	doc := make(map[int]bool, num_terms)
	for i := 0; i < num_terms; i++ {
		doc[rand.Int()] = true
	}

	mh := NewMinhash(256)
	sigs := mh.hash(doc)
	for term, sig := range sigs {
		if sig < 0 {
			t.Errorf("hash has negative signature %q for term %q", sig, term)
		}
	}
}
