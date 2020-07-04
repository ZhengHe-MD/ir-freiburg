package index

import (
	"bufio"
	"math"
	"os"
	"regexp"
	"strings"
)

var (
	nonAlphaCharRegex = regexp.MustCompile("[^a-zA-Z]+")
)

// thread safety is not guaranteed
type InvertedIndex struct {
	invertedLists map[string][]int64
	docs          map[int64]string
}

func NewInvertedIndex() *InvertedIndex {
	return &InvertedIndex{
		invertedLists: make(map[string][]int64),
		docs:          make(map[int64]string),
	}
}

func (ii *InvertedIndex) GetInvertedLists() map[string][]int64 {
	return ii.invertedLists
}

func (ii *InvertedIndex) GetDocByID(id int64) string {
	return ii.docs[id]
}

func (ii *InvertedIndex) ReadFromFile(filename string) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	invertedList := make(map[string][]int64)
	docs := make(map[int64]string)

	scanner := bufio.NewScanner(f)
	docID := int64(0)
	for scanner.Scan() {
		docID += 1

		line := scanner.Text()
		docs[docID] = line
		words := nonAlphaCharRegex.Split(line, -1)
		for _, word := range words {
			word = strings.ToLower(word)
			if _, ok := invertedList[word]; !ok {
				invertedList[word] = nil
			}
			invertedList[word] = append(invertedList[word], docID)
		}
	}

	ii.invertedLists = invertedList
	ii.docs = docs
	return
}

func (ii *InvertedIndex) ProcessQuery(query string) (docIDList []int64, err error) {
	words := nonAlphaCharRegex.Split(query, -1)
	lists := make([][]int64, len(words))
	for i, word := range words {
		lists[i] = ii.invertedLists[word]
	}
	docIDList = MultiMerge(lists...)
	return
}

// MultiMerge merge an arbitrary number of sorted lists.
// TODO: use priority queue to find minimum element.
func MultiMerge(lists ...[]int64) (ret []int64) {
	var idxList = make([]int, len(lists))
	var upperBoundList = make([]int, len(lists))

	for i, list := range lists {
		if len(list) == 0 {
			return
		}
		upperBoundList[i] = len(list)
	}

	checkInBound := func(i int) bool {
		return idxList[i] < upperBoundList[i]
	}

	var lastAdvancedIdx int
	for checkInBound(lastAdvancedIdx) {
		var prevVal, minVal int64 = 0, math.MaxInt64
		var minIdx = 0
		var equal = true
		for i := 0; i < len(lists); i++ {
			val := lists[i][idxList[i]]
			if val < minVal {
				minIdx = i
				minVal = val
			}
			if i > 0 && equal {
				equal = val == prevVal
			}
			prevVal = val
		}

		if equal {
			ret = append(ret, minVal)
		}
		lastAdvancedIdx = minIdx
		idxList[lastAdvancedIdx] += 1
	}

	return
}
