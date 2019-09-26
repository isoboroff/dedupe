package lib

import (
	"log"
	"math"
	"strings"
	"fmt"
)

type LSH struct {
	maps []map[string][]string
	num_rows int
	num_bands int
	num_hashes int
}

func MakeLSH(n, b int) (result LSH) {
	if (n % b) != 0 {
		log.Panicf("Bands %d must divide num_hashes %d evenly\n", b, n)
	}

	result.num_rows = n / b
	result.num_hashes = n
	result.num_bands = b
	result.maps = make([]map[string][]string, result.num_bands)
	for i := range result.maps {
		result.maps[i] = make(map[string][]string)
	}
	log.Printf("LSH with %d hash buckets and %d bands\n", n, b)
	log.Printf("Target Jaccard threshold: %.3f\n",
		math.Pow(1.0 / float64(n), 1.0 / (float64(n) / float64(b))))
	return result
}
	
func MakeLSHForThreshold(n int, thresh float64) LSH {
	num_bands := computeBands(n, thresh)
	return MakeLSH(n, num_bands)
}

func (h LSH) bandprints(hashes []uint32) (prints []string) {
	prints = make([]string, h.num_bands)
	for b := 0; b < h.num_bands; b++ {
		var s strings.Builder
		for r := 0; r < h.num_rows; r++ {
			fmt.Fprintf(&s, "%x", hashes[b * h.num_rows + r])
		}
		prints[b] = s.String()
	}
	return prints
}

func (h LSH) Insert(key string, hashes []uint32) {
	prints := h.bandprints(hashes)
	for b, p := range prints {
		h.maps[b][p] = append(h.maps[b][p], key)
	}
}

func (h LSH) Query(hashes []uint32) []string {
	candidates := make(map[string]bool)
	prints := h.bandprints(hashes)
	for b, p := range prints {
		for _, c := range h.maps[b][p] {
			candidates[c] = true
		}
	}
	result := make([]string, len(candidates))
	i := 0
	for c, _ := range candidates {
		result[i] = c
		i++
	}
	return result
}


func computeBands(num_hashes int, thresh float64) (num_bands int) {
	var t, last_t float64
	for num_bands = num_hashes; num_bands > 1; num_bands-- {
		if ((num_hashes % num_bands) == 0) {
			x := 1.0 / float64(num_bands)
			y := float64(num_bands) / float64(num_hashes)
			last_t = t
			t = math.Pow(x, y)
			// log.Printf("With %d hashes, %d bands gives a threshold of %.3f\n", num_hashes, num_bands, t)
			if t > thresh {
				break
			}
		}
	}
	if !(t >= thresh && last_t <= thresh) {
		log.Printf("Target threshold %.3f not within range [%.3f, %.3f]", thresh, last_t, t)
	}
	return num_bands
}

