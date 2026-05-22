package collection

import (
	"reflect"
	"testing"
)

func TestChunk(t *testing.T) {
	c := Collect([]int{1, 2, 3, 4, 5, 6, 7})
	chunks := c.Chunk(3)

	if len(chunks) != 3 {
		t.Fatalf("Expected 3 chunks, got %d", len(chunks))
	}

	if !reflect.DeepEqual(chunks[0].All(), []int{1, 2, 3}) {
		t.Errorf("Chunk 0 failed")
	}
	if !reflect.DeepEqual(chunks[1].All(), []int{4, 5, 6}) {
		t.Errorf("Chunk 1 failed")
	}
	if !reflect.DeepEqual(chunks[2].All(), []int{7}) {
		t.Errorf("Chunk 2 failed")
	}
}
