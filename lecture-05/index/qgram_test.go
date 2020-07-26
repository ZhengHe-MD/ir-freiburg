package index

import (
	"fmt"
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

func TestQGramIndex_FindMatches(t *testing.T) {
	tests := []struct {
		givenQ                 int
		givenFilename          string
		givenX                 string
		givenDelta             int
		wantMatches            []EntityPEDPair
		wantNumPEDComputations int
	}{
		{
			3, "example.tsv", "frei", 0,
			[]EntityPEDPair{
				{
					Entity: Entity{"frei", 3, "a word"},
					PED:    0,
				},
			},
			1,
		},
		{
			3, "example.tsv", "frei", 2,
			[]EntityPEDPair{
				{
					Entity: Entity{"frei", 3, "a word"},
					PED:    0,
				},
				{
					Entity: Entity{"brei", 2, "another word"},
					PED:    1,
				},
			},
			2,
		},
		{
			3, "example.tsv", "freibu", 2,
			[]EntityPEDPair{
				{
					Entity: Entity{"frei", 3, "a word"},
					PED:    2,
				},
			},
			2,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i+1), func(t *testing.T) {
			q := NewQGramIndex(tt.givenQ)
			assert.NoError(t, q.BuildFromFile(tt.givenFilename))
			matches, numPEDComputations := q.FindMatches(tt.givenX, tt.givenDelta)
			assert.Equal(t, tt.wantMatches, matches, "matches")
			assert.Equal(t, tt.wantNumPEDComputations, numPEDComputations, "numPEDComputations")
		})
	}
}

func TestRankMatches(t *testing.T) {
	tests := []struct {
		givenMatches []EntityPEDPair
		wantMatches  []EntityPEDPair
	}{
		{
			[]EntityPEDPair{
				{Entity: Entity{"foo", 3, "word 0"}, PED: 2},
				{Entity: Entity{"bar", 7, "word 1"}, PED: 0},
				{Entity: Entity{"baz", 2, "word 2"}, PED: 1},
				{Entity: Entity{"boo", 5, "word 3"}, PED: 1},
			},
			[]EntityPEDPair{
				{Entity: Entity{"bar", 7, "word 1"}, PED: 0},
				{Entity: Entity{"boo", 5, "word 3"}, PED: 1},
				{Entity: Entity{"baz", 2, "word 2"}, PED: 1},
				{Entity: Entity{"foo", 3, "word 0"}, PED: 2},
			},
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.wantMatches, RankMatches(tt.givenMatches))
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		givenStr string
		wantStr  string
	}{
		{"Frei, burg !!", "freiburg"},
		{"freiburg", "freiburg"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.wantStr, Normalize(tt.givenStr))
	}
}

func TestQGramIndex_MergeLists(t *testing.T) {
	tests := []struct {
		givenLists [][]int
		wantRet    []WordIdCountPair
	}{
		{
			[][]int{
				{1, 1, 3, 5},
				{2, 3, 3, 9, 9},
			},
			[]WordIdCountPair{
				{1, 2},
				{2, 1},
				{3, 3},
				{5, 1},
				{9, 2},
			},
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.wantRet, MergeLists(tt.givenLists))
	}
}

func TestPrefixEditDistance(t *testing.T) {
	tests := []struct {
		givenX     string
		givenY     string
		givenDelta int
		wantPED    int
	}{
		{"frei", "frei", 0, 0},
		{"frei", "freiburg", 0, 0},
		{"frei", "breifurg", 1, 1},
		{"freiburg", "stuttgart", 2, 3},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.wantPED, PrefixEditDistance(tt.givenX, tt.givenY, tt.givenDelta))
	}
}
