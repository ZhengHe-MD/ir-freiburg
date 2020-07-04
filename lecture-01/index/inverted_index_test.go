package index

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInvertedIndex_ReIndexFromFile(t *testing.T) {
	tests := []struct {
		givenFilename     string
		wantInvertedLists map[string][]int64
	}{
		{
			"../example.txt",
			map[string][]int64{
				"document": {1, 2, 3},
				"first":    {1},
				"second":   {2},
				"third":    {3},
			},
		},
	}

	for _, tt := range tests {
		ii := NewInvertedIndex()
		assert.NoError(t, ii.ReadFromFile(tt.givenFilename))
		assert.Equal(t, tt.wantInvertedLists, ii.invertedLists)
	}
}

func TestMultiMerge(t *testing.T) {
	tests := []struct {
		givenLists [][]int64
		wantRet    []int64
	}{
		{
			[][]int64{
				{1, 3},
				{2, 3},
			},
			[]int64{3},
		},
		{
			[][]int64{
				{1, 3, 4, 6, 7},
				{2, 4, 5, 7, 10},
				{3, 4, 6, 8, 10},
			},
			[]int64{4},
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
		wantDocIDList []int64
		wantErr       error
	}{
		{"../example.txt", "document", []int64{1, 2, 3}, nil},
		{"../example.txt", "document third", []int64{3}, nil},
	}

	for _, tt := range tests {
		ii := NewInvertedIndex()
		assert.NoError(t, ii.ReadFromFile(tt.givenFilename))
		docIDList, err := ii.ProcessQuery(tt.givenQuery)
		assert.NoError(t, err)
		assert.Equal(t, tt.wantDocIDList, docIDList)
	}
}
