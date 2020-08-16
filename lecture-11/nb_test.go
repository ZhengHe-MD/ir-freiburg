package main

import (
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
	"testing"
)

func TestGenerateVocabularies(t *testing.T) {
	tests := []struct {
		givenFilename       string
		wantWordVocabulary  map[string]int
		wantClassVocabulary map[string]int
		wantErr             error
	}{
		{
			"./data/example_train.tsv",
			map[string]int{"a": 0, "b": 1},
			map[string]int{"A": 0, "B": 1},
			nil,
		},
		{
			"./data/example_test.tsv",
			map[string]int{"b": 0, "a": 1},
			map[string]int{"A": 0, "B": 1},
			nil,
		},
	}

	for _, tt := range tests {
		wordVocab, classVocab, err := GenerateVocabularies(tt.givenFilename)
		assert.Equal(t, tt.wantWordVocabulary, wordVocab)
		assert.Equal(t, tt.wantClassVocabulary, classVocab)
		assert.Equal(t, tt.wantErr, err)
	}
}

func TestReadLabeledData(t *testing.T) {
	tests := []struct {
		givenFilename        string
		givenClassVocabulary map[string]int
		givenWordVocabulary  map[string]int
		wantX                *mat.Dense
		wantY                []int
		wantErr              error
	}{
		{
			"./data/example_train.tsv",
			map[string]int{"A": 0, "B": 1},
			map[string]int{"a": 0, "b": 1},
			mat.NewDense(6, 2, []float64{
				2, 1,
				5, 2,
				3, 5,
				3, 2,
				1, 3,
				2, 4,
			}),
			[]int{0, 0, 1, 0, 1, 1},
			nil,
		},
		{
			"./data/example_test.tsv",
			map[string]int{"A": 0, "B": 1},
			map[string]int{"a": 0, "b": 1},
			mat.NewDense(3, 2, []float64{
				6, 1,
				1, 6,
				3, 1,
			}),
			[]int{0, 1, 1},
			nil,
		},
	}

	for _, tt := range tests {
		X, y, err := ReadLabeledData(tt.givenFilename, tt.givenClassVocabulary, tt.givenWordVocabulary)
		assert.Equal(t, tt.wantErr, err)
		assert.Equal(t, tt.wantX, X.ToDense())
		assert.Equal(t, tt.wantY, y)
	}
}

func TestNaiveBayes_Train(t *testing.T) {
	tests := []struct {
		givenTrainFilename string
		wantPWC            []float64
		wantPC             []float64
	}{
		{
			"./data/example_train.tsv",
			[]float64{0.664, 0.336, 0.335, 0.665},
			[]float64{0.500, 0.500},
		},
	}

	for _, tt := range tests {
		wv, cv, err := GenerateVocabularies(tt.givenTrainFilename)
		assert.NoError(t, err)
		X, y, err := ReadLabeledData(tt.givenTrainFilename, cv, wv)
		assert.NoError(t, err)
		nb := &NaiveBayes{}
		assert.NoError(t, nb.Train(X, y, cv, wv))
		assert.InEpsilonSlice(t, tt.wantPWC, nb.PWC, 0.01)
		assert.InEpsilonSlice(t, tt.wantPC, nb.PC, 0.01)
	}
}

func TestNaiveBayes_Predict(t *testing.T) {
	tests := []struct {
		givenTrainFilename string
		givenTestFilename  string
		wantYt             []int
	}{
		{
			"./data/example_train.tsv",
			"./data/example_test.tsv",
			[]int{0, 1, 0},
		},
		{
			"./data/example_train.tsv",
			"./data/example_train.tsv",
			[]int{0, 0, 1, 0, 1, 1},
		},
	}

	for _, tt := range tests {
		wv, cv, err := GenerateVocabularies(tt.givenTrainFilename)
		assert.NoError(t, err)
		X, y, err := ReadLabeledData(tt.givenTrainFilename, cv, wv)
		assert.NoError(t, err)
		nb := &NaiveBayes{}
		assert.NoError(t, nb.Train(X, y, cv, wv))
		Xt, _, err := ReadLabeledData(tt.givenTestFilename, cv, wv)
		assert.NoError(t, err)
		yt, err := nb.Predict(Xt)
		assert.NoError(t, err)
		assert.Equal(t, tt.wantYt, yt)
	}
}

func TestNaiveBayes_Evaluate(t *testing.T) {
	tests := []struct {
		givenTrainFilename string
		givenTestFilename  string
		wantPrecisionList  []float64
		wantRecallList     []float64
		wantF1ScoreList    []float64
	}{
		{
			"./data/example_train.tsv",
			"./data/example_test.tsv",
			[]float64{0.5, 1.0},
			[]float64{1.0, 0.5},
			[]float64{0.67, 0.67},
		},
	}

	for _, tt := range tests {
		wv, cv, err := GenerateVocabularies(tt.givenTrainFilename)
		assert.NoError(t, err)
		X, y, err := ReadLabeledData(tt.givenTrainFilename, cv, wv)
		assert.NoError(t, err)
		Xt, yt, err := ReadLabeledData(tt.givenTestFilename, cv, wv)
		assert.NoError(t, err)
		nb := &NaiveBayes{}
		assert.NoError(t, nb.Train(X, y, cv, wv))
		precisionList, recallList, f1List, err := nb.Evaluate(Xt, yt)
		assert.InEpsilonSlice(t, tt.wantPrecisionList, precisionList, 0.01)
		assert.InEpsilonSlice(t, tt.wantRecallList, recallList, 0.01)
		assert.InEpsilonSlice(t, tt.wantF1ScoreList, f1List, 0.01)
	}
}

func TestLoadStopWordSet(t *testing.T) {
	tests := []struct {
		givenStopWordFilename string
		wantSize              int
		wantSubset            []string
	}{
		{
			"./data/stopwords.txt",
			106,
			[]string{"what", "when", "where"},
		},
	}

	for _, tt := range tests {
		stopWordSet, err := LoadStopWordSet(tt.givenStopWordFilename)
		assert.NoError(t, err)
		assert.Equal(t, tt.wantSize, len(stopWordSet))
		for _, stopWord := range tt.wantSubset {
			assert.True(t, stopWordSet[stopWord])
		}
	}
}
