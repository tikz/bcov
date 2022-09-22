package cov

import (
	"testing"
)

func TestExonDepth(t *testing.T) {
	exon := ExonDepth("../tests/little.bam", "1", 100, 100000)
	testPos := func(pos int, expectedCount int) {
		count := exon.PositionDepth[Position(pos)]
		if count != Depth(expectedCount) {
			t.Errorf("Expected pos %d read count to be %d, got %d", pos, expectedCount, count)
		}
	}

	testPos(70017, 16)
	testPos(69956, 40)
	testPos(99270, 1)
}
