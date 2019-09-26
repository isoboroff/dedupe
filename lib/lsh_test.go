
package lib

import (
	"testing"
	"math"
)

func TestCompute(t *testing.T) {
	target := 0.90
	h := MakeLSHForThreshold(256, target)

	next_band := h.num_bands
	for next_band++ ; h.num_hashes % next_band != 0; next_band++ {
	}
	
	x0 := 1.0 / float64(next_band)
	y0 := float64(next_band) / float64(h.num_hashes)
	low_threshold := math.Pow(x0, y0)

	if low_threshold > target {
		t.Errorf("Low threshold of %.3f is too high to hit target %.3f",
			low_threshold, target)
	}
	
	x := 1.0 / float64(h.num_bands)
	y := float64(h.num_bands) / float64(h.num_hashes)
	threshold := math.Pow(x, y)

	if threshold < low_threshold {
		t.Errorf("low_threshold is %.3f, threshold is %.3f",
			low_threshold, threshold)
	}
}
