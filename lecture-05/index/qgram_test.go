package index

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewQGramIndex(t *testing.T) {
	tests := []struct {
		givenQ            int
		wantQ             int
		wantPadding       string
		wantInvertedLists map[string][]int
	}{
		{3, 3, "$$", make(map[string][]int)},
	}

	for _, tt := range tests {
		q := NewQGramIndex(tt.givenQ)
		assert.Equal(t, tt.wantQ, q.Q)
		assert.Equal(t, tt.wantPadding, q.Padding)
		assert.Equal(t, tt.wantInvertedLists, q.InvertedLists)
	}
}

func TestQGramIndex_BuildFromFile(t *testing.T) {
	tests := []struct {
		givenFilename    string
		givenQ           int
		wantInvertedList map[string][]int
	}{
		{
			"example.tsv",
			3,
			map[string][]int{
				"$$b": {2},
				"$$f": {1},
				"$br": {2},
				"$fr": {1},
				"bre": {2},
				"ei$": {1, 2},
				"fre": {1},
				"i$$": {1, 2},
				"rei": {1, 2},
			},
		},
	}

	for _, tt := range tests {
		q := NewQGramIndex(tt.givenQ)
		err := q.BuildFromFile(tt.givenFilename)
		assert.NoError(t, err)
		assert.Equal(t, tt.wantInvertedList, q.InvertedLists)
	}
}

func TestQGramIndex_ComputeQGram(t *testing.T) {
	tests := []struct {
		givenQ        int
		givenWord     string
		wantQGramList []string
	}{
		{
			3, "freiburg",
			[]string{
				"$$f", "$fr", "fre", "rei", "eib",
				"ibu", "bur", "urg", "rg$", "g$$",
			},
		},
	}

	for _, tt := range tests {
		q := NewQGramIndex(tt.givenQ)
		assert.Equal(t, tt.wantQGramList, q.ComputeQGram(tt.givenWord))
	}
}

func TestQGramIndex_normalize(t *testing.T) {
	tests := []struct{
		givenStr string
		wantStr string
	} {
		{"Frei, burg !!", "freiburg"},
		{"freiburg", "freiburg"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.wantStr, Normalize(tt.givenStr))
	}
}