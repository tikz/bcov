package bed

import (
	"testing"
)

func TestReader(t *testing.T) {
	reader, err := NewReader("../tests/test.bed")
	if err != nil {
		t.Errorf("Error opening bed file: %v", err)
	}

	exon, err := reader.Read()
	if err != nil {
		t.Errorf("Error reading bed file: %v", err)
	}

	if exon.Chromosome != "1" {
		t.Errorf("Expected chromosome to be 1, got %s", exon.Chromosome)
	}

	if exon.Start != 35100 {
		t.Errorf("Expected start to be 35100, got %d", exon.Start)
	}

	if exon.End != 35110 {
		t.Errorf("Expected end to be 35110, got %d", exon.End)
	}

	if exon.Name != "MOCK1" {
		t.Errorf("Expected name to be MOCK1, got %s", exon.Name)
	}
}

func TestOpen(t *testing.T) {
	_, err := NewReader("xxx")
	if err == nil {
		t.Errorf("Expected error opening file, got nil")
	}
}
