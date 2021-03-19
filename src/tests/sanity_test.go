package tests

import (
	"math"
	"testing"
)

// NOTE(Jovan): Sample test
func TestSanityAbs(t *testing.T) {
	t.Log("Running SanityAbs test...")
	got := math.Abs(-1)
	if got != 1 {
		t.Errorf("Abs(-1) = %f; want 1", got)
	}
}
