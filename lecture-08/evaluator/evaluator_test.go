package evaluator

import (
	"github.com/ZhengHe-MD/ir-freiburg.git/lecture-08/index"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

const epsilon = 0.001

func TestReadBenchmark(t *testing.T) {
	tests := []struct {
		givenFilename string
		wantBenchmark map[string]map[int64]interface{}
	}{
		{
			"example-benchmark.txt",
			map[string]map[int64]interface{}{
				"animated film": {
					1: struct{}{},
					3: struct{}{},
					4: struct{}{},
				},
				"short film": {
					3: struct{}{},
					4: struct{}{},
				},
			},
		},
	}

	for _, tt := range tests {
		queryDocs, err := ReadBenchmark(tt.givenFilename)
		assert.NoError(t, err)
		assert.Equal(t, tt.wantBenchmark, queryDocs)
	}
}

func TestPrecisionAtK(t *testing.T) {
	tests := []struct {
		givenResultIds   []int64
		givenRelevantIds map[int64]interface{}
		givenK           int
		wantRetPrecision float64
	}{
		{
			[]int64{5, 3, 6, 1, 2},
			map[int64]interface{}{
				1: struct{}{},
				2: struct{}{},
				5: struct{}{},
				6: struct{}{},
				7: struct{}{},
				8: struct{}{},
			},
			2,
			0.5,
		},
		{
			[]int64{5, 3, 6, 1, 2},
			map[int64]interface{}{
				1: struct{}{},
				2: struct{}{},
				5: struct{}{},
				6: struct{}{},
				7: struct{}{},
				8: struct{}{},
			},
			4,
			0.75,
		},
		{
			[]int64{},
			map[int64]interface{}{
				1: struct{}{},
				2: struct{}{},
				5: struct{}{},
			},
			4,
			0.0,
		},
		{
			[]int64{5, 3, 6, 1, 2},
			map[int64]interface{}{
				1: struct{}{},
				2: struct{}{},
				5: struct{}{},
				6: struct{}{},
				7: struct{}{},
				8: struct{}{},
			},
			0,
			0.0,
		},
	}

	for _, tt := range tests {
		p := PrecisionAtK(tt.givenResultIds, tt.givenRelevantIds, tt.givenK)
		assert.Equal(t, tt.wantRetPrecision, p)
	}
}

func TestAveragePrecision(t *testing.T) {
	tests := []struct {
		givenResultIds   []int64
		givenRelevantIds map[int64]interface{}
		wantAP           float64
	}{
		{
			[]int64{7, 17, 9, 42, 5},
			map[int64]interface{}{
				5:  struct{}{},
				7:  struct{}{},
				12: struct{}{},
				42: struct{}{},
			},
			0.525,
		},
	}

	for _, tt := range tests {
		ap := AveragePrecision(tt.givenResultIds, tt.givenRelevantIds)
		assert.Equal(t, tt.wantAP, ap)
	}
}

func TestEvaluate(t *testing.T) {
	tests := []struct {
		givenFilename          string
		givenBenchmarkFilename string
		givenBM25B             float64
		givenBM25K             float64
		wantMPS                MPS
	}{
		{
			"example.txt", "example-benchmark.txt",
			0.75, 1.25,
			MPS{0.667, 0.833, 0.694},
		},
	}

	for _, tt := range tests {
		ii := index.NewInvertedIndex()
		err := ii.ReadFromFile(tt.givenFilename, tt.givenBM25B, tt.givenBM25K, index.RefinementOptions{})
		assert.NoError(t, err)

		benchmark, err := ReadBenchmark(tt.givenBenchmarkFilename)
		assert.NoError(t, err)

		mps := Evaluate(ii, benchmark, index.RefinementOptions{})
		assert.True(t, math.Abs(tt.wantMPS.MPAt3-mps.MPAt3) <= epsilon)
		assert.True(t, math.Abs(tt.wantMPS.MPAtR-mps.MPAtR) <= epsilon)
		assert.True(t, math.Abs(tt.wantMPS.MAP-mps.MAP) <= epsilon)
	}
}
