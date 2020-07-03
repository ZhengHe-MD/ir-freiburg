package index

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestInvertedIndex_ReIndexFromFile(t *testing.T) {
	tests := []struct {
		givenFilename     string
		givenBM25B        float64
		givenBM25K        float64
		wantInvertedLists map[string][]Posting
	}{
		{
			"../example.txt",
			0, math.Inf(1),
			map[string][]Posting{
				"animated": {
					{1, 0.415},
					{2, 0.415},
					{4, 0.415},
				},
				"animation": {
					{3, 2.0},
				},
				"film": {
					{2, 1.0},
					{4, 1.0},
				},
				"movie": {
					{1, 0.0},
					{2, 0.0},
					{3, 0.0},
					{4, 0.0},
				},
				"non": {
					{2, 2.0},
				},
				"short": {
					{3, 1.0},
					{4, 2.0},
				},
			},
		},
		{
			"../example.txt",
			0.75, 1.75,
			map[string][]Posting{
				"animated": {
					{1, 0.459},
					{2, 0.402},
					{4, 0.358},
				},
				"animation": {
					{3, 2.211},
				},
				"film": {
					{2, 0.969},
					{4, 0.863},
				},
				"movie": {
					{1, 0.0},
					{2, 0.0},
					{3, 0.0},
					{4, 0.0},
				},
				"non": {
					{2, 1.938},
				},
				"short": {
					{3, 1.106},
					{4, 1.313},
				},
			},
		},
	}

	for _, tt := range tests {
		ii := NewInvertedIndex()
		assert.NoError(t, ii.ReIndexFromFile(tt.givenFilename, tt.givenBM25B, tt.givenBM25K))
		assert.Equal(t, tt.wantInvertedLists, ii.getRoundedInvertedIndex())
	}
}

func TestMerge(t *testing.T) {
	tests := []struct {
		givenPostingsA []Posting
		givenPostingsB []Posting
		wantPostings   []Posting
	}{
		{
			[]Posting{
				{1, 2.1},
				{5, 3.2},
			},
			[]Posting{
				{1, 1.7},
				{2, 1.3},
				{5, 3.3},
			},
			[]Posting{
				{1, 3.8},
				{2, 1.3},
				{5, 6.5},
			},
		},
	}

	for _, tt := range tests {
		postings := Merge(tt.givenPostingsA, tt.givenPostingsB)
		assert.Equal(t, tt.wantPostings, postings)
	}
}

func TestInvertedIndex_ProcessQuery(t *testing.T) {
	tests := []struct {
		givenFilename string
		givenQuery    string
		givenBM25B    float64
		givenBM25K    float64
		wantDocIDList []int
		wantErr       error
	}{
		{"../example.txt", "movie short", 0.75, 1.75, []int{4, 3, 2, 1}, nil},
		{"../example.txt", "non animated film", 0.75, 1.75, []int{2, 4, 1}, nil},
	}

	for _, tt := range tests {
		ii := NewInvertedIndex()
		assert.NoError(t, ii.ReIndexFromFile(tt.givenFilename, tt.givenBM25B, tt.givenBM25K))
		docIDList, err := ii.ProcessQuery(tt.givenQuery)
		assert.NoError(t, err)
		assert.Equal(t, tt.wantDocIDList, docIDList)
	}
}
