package postinglist

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPostingList_ReadFromFile(t *testing.T) {
	tests := []struct{
		givenFilename string
		wantPostingList PostingList
	} {
		{
			"../data/example1.txt",
			PostingList{
				cap:       3,
				num:       3,
				docIDList: []int64{2, 3, 6},
				scoreList: []int64{5, 1, 2},
			},
		},
		{
			"../data/example2.txt",
			PostingList{
				cap:       4,
				num:       4,
				docIDList: []int64{1, 2, 4, 6},
				scoreList: []int64{1, 4, 3, 3},
			},
		},
		{
			"../data/example3.txt",
			PostingList{
				cap:       2,
				num:       2,
				docIDList: []int64{5, 7},
				scoreList: []int64{1, 2},
			},
		},
	}

	for _, tt := range tests {
		p := NewPostingList()
		err := p.ReadFromFile(tt.givenFilename)
		assert.NoError(t, err)
		assert.Equal(t, tt.wantPostingList, *p)
	}
}