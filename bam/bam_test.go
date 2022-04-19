package bam

import (
	"testing"
)

func TestReader(t *testing.T) {
	reader, err := NewReader("../tests/little.bam")
	if err != nil {
		t.Errorf("Error opening bed file: %v", err)
	}

	record, err := reader.Read()
	if err != nil {
		t.Errorf("Error reading bed file: %v", err)
	}

	if record.Name != "A00564:437:H7W3JDSX3:1:2466:4933:31172" {
		t.Errorf("Expected #1 record name to be A00564:437:H7W3JDSX3:1:2466:4933:31172, got %s", record.Name)
	}

	for i := 0; i < 1000; i++ {
		record, err = reader.Read()
	}

	if record.Name != "A00564:437:H7W3JDSX3:1:2465:1353:31767" {
		t.Errorf("Expected #1001 record name to be A00564:437:H7W3JDSX3:1:2465:1353:31767, got %s", record.Name)
	}

}
