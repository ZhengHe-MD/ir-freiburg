package evaluator

import (
	"bufio"
	"github.com/ZhengHe-MD/ir-freiburg.git/lecture-02/index"
	"log"
	"os"
	"strconv"
	"strings"
)

type MPS struct {
	MPAt3 float64
	MPAtR float64
	MAP   float64
}

func Evaluate(ii *index.InvertedIndex, benchmark map[string]map[int64]interface{}, useRefinements bool) (mps MPS) {
	var PAt3SoFar, PAtRSoFar, APSoFar float64

	for query, relevantIds := range benchmark {
		postings := ii.ProcessQuery(query)
		var resultIds []int64
		for _, posting := range postings {
			resultIds = append(resultIds, posting.DocID)
		}

		PAt3SoFar += PrecisionAtK(resultIds, relevantIds, 3)
		PAtRSoFar += PrecisionAtK(resultIds, relevantIds, len(relevantIds))
		APSoFar += AveragePrecision(resultIds, relevantIds)
	}

	queryNum := float64(len(benchmark))

	mps = MPS{
		MPAt3: PAt3SoFar / queryNum,
		MPAtR: PAtRSoFar / queryNum,
		MAP:   APSoFar / queryNum,
	}
	return
}

func PrecisionAtK(resultIds []int64, relevantIds map[int64]interface{}, k int) (pk float64) {
	if len(resultIds) == 0 || k == 0 {
		return
	}

	ub := k
	if ub > len(resultIds) {
		ub = len(resultIds)
	}

	var relevantNum float64
	for _, resultId := range resultIds[:ub] {
		if _, ok := relevantIds[resultId]; ok {
			relevantNum += 1.0
		}
	}

	pk = relevantNum / float64(ub)
	return
}

func AveragePrecision(resultIds []int64, relevantIds map[int64]interface{}) (ap float64) {
	var pAtRSum float64
	var relNumSoFar, numSoFar float64

	for _, resultId := range resultIds {
		numSoFar += 1
		if _, ok := relevantIds[resultId]; ok {
			relNumSoFar += 1.0
			pAtRSum += relNumSoFar / numSoFar
		}
	}

	ap = pAtRSum / float64(len(relevantIds))
	return
}

func ReadBenchmark(filename string) (benchmark map[string]map[int64]interface{}, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	benchmark = make(map[string]map[int64]interface{})

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, "\t")
		if len(parts) != 2 {
			log.Printf("invalid line %s\n", line)
			continue
		}

		query := parts[0]

		var docId int64
		for _, docIdStr := range strings.Fields(parts[1]) {
			docId, err = strconv.ParseInt(docIdStr, 10, 64)
			if err != nil {
				return
			}
			if benchmark[query] == nil {
				benchmark[query] = make(map[int64]interface{})
			}
			benchmark[query][docId] = struct{}{}
		}
	}
	return
}
