package index

import (
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
	"math"
	"testing"
)

const epsilon = 1E-04

func TestInvertedIndex_ReadFromFile(t *testing.T) {
	tests := []struct {
		givenFilename     string
		givenBM25B        float64
		givenBM25K        float64
		wantInvertedLists map[string][]Posting
		wantTerms         []string
		wantNumTerm       int
		wantNumDocs       int
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
			[]string{"movie", "animated", "non", "film", "short", "animation"}, 6, 4,
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
			[]string{"movie", "animated", "non", "film", "short", "animation"}, 6, 4,
		},
	}

	for _, tt := range tests {
		ii := NewInvertedIndex()
		assert.NoError(t, ii.ReadFromFile(tt.givenFilename, tt.givenBM25B, tt.givenBM25K, RefinementOptions{RankingScore: RankingScoreBM25}))
		assert.Equal(t, tt.wantInvertedLists, ii.getRoundedInvertedIndex())
		assert.Equal(t, tt.wantTerms, ii.terms)
		assert.Equal(t, tt.wantNumTerm, ii.numTerms)
		assert.Equal(t, tt.wantNumDocs, ii.numDocs)
	}
}

func TestInvertedIndex_PreprocessingVSM(t *testing.T) {
	tests := []struct {
		givenFilename      string
		givenBM25B         float64
		givenBM25K         float64
		givenNormalization Normalization
		wantMatrix         *mat.Dense
	}{
		{
			"example.txt",
			0, math.Inf(1), None,
			mat.NewDense(6, 4, []float64{
				0.000, 0.000, 0.000, 0.000,
				0.415, 0.415, 0.000, 0.415,
				0.000, 2.000, 0.000, 0.000,
				0.000, 1.000, 0.000, 1.000,
				0.000, 0.000, 1.000, 2.000,
				0.000, 0.000, 2.000, 0.000,
			}),
		},
		{
			"example.txt",
			0.75, 1.75, ColumnWiseL2,
			mat.NewDense(6, 4, []float64{
				0.000, 0.000, 0.000, 0.000,
				1.000, 0.182, 0.000, 0.222,
				0.000, 0.879, 0.000, 0.000,
				0.000, 0.440, 0.000, 0.535,
				0.000, 0.000, 0.447, 0.815,
				0.000, 0.000, 0.894, 0.000,
			}),
		},
	}

	for _, tt := range tests {
		ii := NewInvertedIndex()
		assert.NoError(t, ii.ReadFromFile(tt.givenFilename, tt.givenBM25B, tt.givenBM25K, RefinementOptions{
			RankingScore: RankingScoreBM25,
		}))
		ii.PreprocessingVSM(tt.givenNormalization)
		assert.Equal(t, tt.wantMatrix, ii.getRoundedTDMatrix().ToDense())
	}
}

func TestInvertedIndex_ProcessQueryVSM(t *testing.T) {
	tests := []struct {
		givenQuery         string
		givenInvertedLists map[string][]Posting
		givenNumTerms      int
		givenNumDocs       int
		givenTerms         []string
		givenTermToIdx     map[string]int
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
			3, 3,
			[]string{"foo", "bar", "baz"}, map[string]int{"foo": 0, "bar": 1, "baz": 2},
			[]Posting{
				{3, 1.1},
				{2, 0.7},
				{1, 0.6},
			},
		},
		{
			"foo bar foo bar",
			map[string][]Posting{
				"foo": {
					{1, 0.2},
					{3, 0.6},
				},
				"bar": {
					{2, 0.4},
					{3, 0.1},
					{4, 0.8},
				},
			},
			2, 4,
			[]string{"foo", "bar"}, map[string]int{"foo": 0, "bar": 1},
			[]Posting{
				{4, 1.6},
				{3, 1.4},
				{2, 0.8},
				{1, 0.4},
			},
		},
	}

	for _, tt := range tests {
		ii := InvertedIndex{
			invertedLists: tt.givenInvertedLists,
			numTerms:      tt.givenNumTerms,
			numDocs:       tt.givenNumDocs,
			terms:         tt.givenTerms,
			termToIdx:     tt.givenTermToIdx,
		}
		ii.PreprocessingVSM(None)
		docPostings := ii.ProcessQueryVSM(tt.givenQuery, RefinementOptions{RankingScore: RankingScoreBM25})
		for i, wantPosting := range tt.wantResultPosting {
			assert.Equal(t, wantPosting.DocID, docPostings[i].DocID)
			assert.True(t, math.Abs(wantPosting.Score-docPostings[i].Score) < epsilon)
		}
	}
}
