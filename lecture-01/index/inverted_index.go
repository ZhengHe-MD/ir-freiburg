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
	invertedLists map[string][]int
	docs          map[int]string
}

func NewInvertedIndex() *InvertedIndex {
	return &InvertedIndex{
		invertedLists: make(map[string][]int),
		docs:          make(map[int]string),
	}
}

func (ii *InvertedIndex) GetInvertedLists() map[string][]int {
	return ii.invertedLists
}

func (ii *InvertedIndex) GetDocByID(id int) string {
	return ii.docs[id]
}

func (ii *InvertedIndex) ReIndexFromFile(filename string) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	invertedList := make(map[string][]int)

	scanner := bufio.NewScanner(f)
	docID := 0
	for scanner.Scan() {
		docID += 1
		line := scanner.Text()
		ii.docs[docID] = line
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
	return
}

func (ii *InvertedIndex) ProcessQuery(query string) (docIDList []int, err error) {
	words := nonAlphaCharRegex.Split(query, -1)
	lists := make([][]int, len(words))
	for i, word := range words {
		lists[i] = ii.invertedLists[word]
	}
	docIDList = MultiMerge(lists...)
	return
}

func MultiMerge(lists ...[]int) (ret []int) {
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
		var prevVal, minVal, minIdx = 0, math.MaxInt32, 0
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
