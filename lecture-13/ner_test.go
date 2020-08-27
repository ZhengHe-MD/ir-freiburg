package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNER_POSTag(t *testing.T) {
	tests := []struct {
		givenTransitionProbabilitiesFilename string
		givenWordDistributionFilename        string
		givenSentence                        []string
		wantWordTags                         []WordTag
	}{
		{
			"./example-trans-probs.tsv",
			"./example-word-distrib.tsv",
			[]string{"James", "Bond", "is", "an", "agent"},
			[]WordTag{
				{"James", "NNP"},
				{"Bond", "NNP"},
				{"is", "VB"},
				{"an", "OTHER"},
				{"agent", "NN"},
			},
		},
	}

	for _, tt := range tests {
		transitionProbs, err := ReadTransitionProbabilitiesFromFile(tt.givenTransitionProbabilitiesFilename)
		assert.NoError(t, err)
		wordDistrib, err := ReadWordDistributionFromFile(tt.givenWordDistributionFilename)
		assert.NoError(t, err)
		ner := NewNER(transitionProbs, wordDistrib)
		wordTags, err := ner.POSTag(tt.givenSentence)
		assert.NoError(t, err)
		assert.Equal(t, tt.wantWordTags, wordTags)
	}
}

func TestNER_FindNamedEntities(t *testing.T) {
	tests := []struct {
		givenTransitionProbabilitiesFilename string
		givenWordDistributionFilename string
		givenSentence []string
		wantNamedEntities []string
	} {
		{
			"./example-trans-probs.tsv",
			"./example-word-distrib.tsv",
			[]string{"James", "Bond", "is", "an", "agent"},
			[]string{"James Bond"},
		},
	}

	for _, tt := range tests {
		transitionProbs, err := ReadTransitionProbabilitiesFromFile(tt.givenTransitionProbabilitiesFilename)
		assert.NoError(t, err)
		wordDistrib, err := ReadWordDistributionFromFile(tt.givenWordDistributionFilename)
		assert.NoError(t, err)
		ner := NewNER(transitionProbs, wordDistrib)
		entities, err := ner.FindNamedEntities(tt.givenSentence)
		assert.NoError(t, err)
		assert.Equal(t, tt.wantNamedEntities, entities)
	}
}

func TestNER_ReadTransitionProbabilitiesFromFile(t *testing.T) {
	tests := []struct {
		givenFilename string
		wantProbs     map[string]map[string]float64
	}{
		{
			"./example-trans-probs.tsv",
			map[string]map[string]float64{
				"NNP":   {"NNP": 0.5, "VB": 0.3, "END": 0.2},
				"BEG":   {"NNP": 0.4, "NN": 0.4, "OTHER": 0.2},
				"NN":    {"NN": 0.4, "VB": 0.4, "END": 0.2},
				"VB":    {"NNP": 0.4, "NN": 0.4, "OTHER": 0.2},
				"OTHER": {"NNP": 0.1, "NN": 0.5, "OTHER": 0.4},
			},
		},
	}

	for _, tt := range tests {
		probs, err := ReadTransitionProbabilitiesFromFile(tt.givenFilename)
		assert.NoError(t, err)
		assert.Equal(t, tt.wantProbs, probs)
	}
}

func TestNER_ReadWordDistributionFromFile(t *testing.T) {
	tests := []struct {
		givenFilename   string
		wantWordDistrib map[string]map[string]float64
	}{
		{
			"./example-word-distrib.tsv",
			map[string]map[string]float64{
				"NNP":   {"James": 0.5, "Bond": 0.5},
				"NN":    {"agent": 0.9, "James": 0.1},
				"VB":    {"is": 0.8, "Bond": 0.2},
				"OTHER": {"James": 0.2, "Bond": 0.2, "an": 0.6},
			},
		},
	}

	for _, tt := range tests {
		wordDistrib, err := ReadWordDistributionFromFile(tt.givenFilename)
		assert.NoError(t, err)
		assert.Equal(t, tt.wantWordDistrib, wordDistrib)
	}
}
