package lib

import (
	"hash"
	"hash/fnv"
	"fmt"
	"strings"	
)

const (
	SHINGLE_LEN = 9
)

var hashfn hash.Hash32

func fingerprint(s string) uint32 {
	hashfn := fnv.New32()
	hashfn.Write([]byte(s))
	return hashfn.Sum32()
}

func Shingle(s string) []uint32 {
	if len(s) == 0 {
		return nil
	}
	f := strings.Fields(s)
	if len(f) == 0 {
		return nil
	}
	num_shingles := len(f) - SHINGLE_LEN + 1
	if num_shingles < 0 {
		num_shingles = 1
	}
	resmap := make(map[uint32]bool, num_shingles)
	for i := 0; i < num_shingles; i++ {
		var shingle string
		if len(f) < SHINGLE_LEN {
			shingle = strings.Join(f, " ")
		} else {
			shingle = strings.Join(f[i:i+SHINGLE_LEN], " ")
		}
		resmap[fingerprint(shingle)] = true
	}
	result := make([]uint32, len(resmap))
	i := 0
	for key, _ := range resmap {
		result[i] = key
		i++
	}
	return result
}

func ShingleChars(s string) []uint32 {
	lens := len(s)
	if lens < SHINGLE_LEN {
		s = fmt.Sprintf(fmt.Sprintf("%%%ds", SHINGLE_LEN), s)
		lens = len(s)
	}
	num_shingles := lens - SHINGLE_LEN
	result := make([]uint32, num_shingles)
	for i := 0; i < num_shingles; i++ {
		shingle := s[i:i+SHINGLE_LEN]
		result[i] = fingerprint(shingle)
	}
	return result
}
