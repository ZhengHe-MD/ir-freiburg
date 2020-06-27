package index

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInvertedIndex_ReIndexFromFile(t *testing.T) {
	tests := []struct {
		givenFilename     string
		wantInvertedLists map[string][]int
	}{
		{
			"../fixtures.txt",
			map[string][]int{
				"document": {1, 2, 3},
				"first":    {1},
				"second":   {2},
				"third":    {3},
			},
		},
	}

	for _, tt := range tests {
		ii := NewInvertedIndex()
		assert.NoError(t, ii.ReIndexFromFile(tt.givenFilename))
		assert.Equal(t, tt.wantInvertedLists, ii.invertedLists)
	}
}

func TestMultiMerge(t *testing.T) {
	tests := []struct {
		givenLists [][]int
		wantRet    []int
	}{
		{
			[][]int{
				{1, 3},
				{2, 3},
			},
			[]int{3},
		},
		{
			[][]int{
				{1, 3, 4, 6, 7},
				{2, 4, 5, 7, 10},
				{3, 4, 6, 8, 10},
			},
			[]int{4},
		},
	}

	for _, tt := range tests {
		ret := MultiMerge(tt.givenLists...)
		assert.Equal(t, tt.wantRet, ret)
	}
}

func TestInvertedIndex_ProcessQuery(t *testing.T) {
	tests := []struct {
		givenFilename string
		givenQuery    string
		wantDocIDList []int
		wantErr       error
	}{
		{"../fixtures.txt", "document", []int{1, 2, 3}, nil},
		{"../fixtures.txt", "document third", []int{3}, nil},
	}

	for _, tt := range tests {
		ii := NewInvertedIndex()
		assert.NoError(t, ii.ReIndexFromFile(tt.givenFilename))
		docIDList, err := ii.ProcessQuery(tt.givenQuery)
		assert.NoError(t, err)
		assert.Equal(t, tt.wantDocIDList, docIDList)
	}
}
