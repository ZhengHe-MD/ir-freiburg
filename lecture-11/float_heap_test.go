package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCappedFloatHeap_All(t *testing.T) {
	h := NewCappedItemHeap(3)

	h.Push(&Item{
		ID:       3,
		Priority: 5,
	})
	h.Push(&Item{
		ID:       1,
		Priority: 10,
	})
	h.Push(&Item{
		ID:       5,
		Priority: 3,
	})
	h.Push(&Item{
		ID:       100,
		Priority: 1,
	})

	item1, ok1 := h.Pop()
	assert.True(t, ok1)
	assert.Equal(t, 5, item1.ID)

	item2, ok2 := h.Pop()
	assert.True(t, ok2)
	assert.Equal(t, 3, item2.ID)

	item3, ok3 := h.Pop()
	assert.True(t, ok3)
	assert.Equal(t, 1, item3.ID)

	_, ok4 := h.Pop()
	assert.False(t, ok4)
}
