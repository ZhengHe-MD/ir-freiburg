package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInvertedIndex_ReindexFromFile(t *testing.T) {
	tests := []struct{
		givenFilename string
		wantInvertedLists map[string][]int
	} {
		{
			"fixtures.txt",
			map[string][]int{
				"document": {1, 2, 3},
				"first": {1},
				"second": {2},
				"third": {3},
			},
		},
	}

	for _, tt := range tests {
		ii := NewInvertedIndex()
		assert.NoError(t, ii.ReIndexFromFile(tt.givenFilename))
		assert.Equal(t, tt.wantInvertedLists, ii.invertedLists)
	}
}
