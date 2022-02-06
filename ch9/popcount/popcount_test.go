package popcount

import "testing"

func TestParallelPopcount(t *testing.T) {
	result := ParallelPopCount(8)
	if result != 1 {
		t.Errorf("expect 1, got %v", result)
	}
}