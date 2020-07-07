package index

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

const epsilon = 1E-04

func TestInvertedIndex_ReIndexFromFile(t *testing.T) {
	tests := []struct {
		givenFilename     string
		givenBM25B        float64
		givenBM25K        float64
		wantInvertedLists map[string][]Posting
	}{
		{
			"example.txt",
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
			"example.txt",
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
		assert.NoError(t, ii.ReadFromFile(tt.givenFilename, tt.givenBM25B, tt.givenBM25K, RefinementOptions{}))
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
		givenQuery         string
		givenInvertedLists map[string][]Posting
		wantResultPosting  []Posting
	}{
		{
			"foo bar",
			map[string][]Posting{
				"foo": {
					{1, 0.2},
					{3, 0.6},
				},
				"bar": {
					{1, 0.4},
					{2, 0.7},
					{3, 0.5},
				},
				"baz": {
					{2, 0.1},
				},
			},
			[]Posting{
				{3, 1.1},
				{2, 0.7},
				{1, 0.6},
			},
		},
	}

	for _, tt := range tests {
		ii := InvertedIndex{invertedLists: tt.givenInvertedLists}
		docPostings := ii.ProcessQuery(tt.givenQuery, RefinementOptions{})
		for i, wantPosting := range tt.wantResultPosting {
			assert.Equal(t, wantPosting.DocID, docPostings[i].DocID)
			assert.True(t, math.Abs(wantPosting.BM25-docPostings[i].BM25) < epsilon)
		}
	}
}

func TestInvertedIndex_Process(t *testing.T) {
	tests := []struct {
		givenFilename   string
		givenQuery      string
		givenBM25B      float64
		givenBM25K      float64
		wantDocPostings []Posting
	}{
		{
			"example.txt", "movie short",
			0.75, 1.75,
			[]Posting{
				{4, 0.0},
				{3, 0.0},
				{2, 0.0},
				{1, 0.0},
			},
		},
		{
			"example.txt",
			"non animated film",
			0.75, 1.75,
			[]Posting{
				{2, 0.0},
				{4, 0.0},
				{1, 0.0},
			},
		},
	}

	for _, tt := range tests {
		ii := NewInvertedIndex()
		assert.NoError(t, ii.ReadFromFile(tt.givenFilename, tt.givenBM25B, tt.givenBM25K, RefinementOptions{}))
		docPostings := ii.ProcessQuery(tt.givenQuery, RefinementOptions{})
		for i, wantDocPosting := range tt.wantDocPostings {
			assert.Equal(t, wantDocPosting.DocID, docPostings[i].DocID)
		}
	}
}
