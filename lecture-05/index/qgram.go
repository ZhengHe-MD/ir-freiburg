package index

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Entity struct {
	Name        string
	Score       int
	Description string
}

type QGramIndex struct {
	Q             int
	Padding       string
	InvertedLists map[string][]int
	EntityMap     map[int]Entity
}

// NewQGramIndex creates an empty QGramIndex.
func NewQGramIndex(q int) *QGramIndex {
	return &QGramIndex{
		Q:             q,
		Padding:       strings.Repeat("$", q-1),
		InvertedLists: make(map[string][]int),
		EntityMap:     make(map[int]Entity),
	}
}

// BuildFromFile builds index from given file (one line per entity, see ES5).
func (q *QGramIndex) BuildFromFile(filename string) (err error) {
	if q == nil {
		err = errors.New("got nil QGramIndex")
		return
	}

	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		err = errors.New("file less than one line")
		return
	}

	// ignore the first line
	_ = scanner.Text()

	var wordId int
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		// fail fast for now
		if len(parts) < 3 {
			err = errors.New("invalid line format")
			return
		}
		wordId += 1

		var word = parts[0]
		for _, qgram := range q.ComputeQGram(word) {
			q.InvertedLists[qgram] = append(q.InvertedLists[qgram], wordId)
		}

		var score int64
		score, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return
		}

		var description = parts[2]
		q.EntityMap[wordId] = Entity{
			Name:        word,
			Score:       int(score),
			Description: description,
		}
	}

	return
}

// ComputeQGram computes q-grams for padded, normalized version of given string.
func (q *QGramIndex) ComputeQGram(word string) (qGramList []string) {
	word = fmt.Sprintf("%s%s%s", q.Padding, Normalize(word), q.Padding)
	for i := 0; i < len(word)-q.Q+1; i++ {
		qGramList = append(qGramList, word[i:i+q.Q])
	}
	return
}

type EntityPEDPair struct {
	Entity Entity
	PED    int
}

// Find all entities y with PED(x, y) ≤ δ for the given string x and a given
// integer δ. First use the q-gram index to exclude all entities that do not
// have a sufficient number of q-grams in common with x, as explained in the
// lecture. Then, compute the PED only for the remaining candidate entities.
// The method should record the number of PED computations as well. Return a
// pair (matches, num_ped_computations), where (1) 'matches' is a list of
// (entity, ped) pairs, where 'entity' is a matching entity y with
// PED(x, y) ≤ δ and 'ped' is the actual PED value; (2) 'num_ped_computations'
// is the number of PED computations done while computing the result.
func (q *QGramIndex) FindMatches(x string, delta int) (matches []EntityPEDPair, numPEDComputations int) {
	var lists [][]int
	for _, qGram := range q.ComputeQGram(x) {
		if invertedList, ok := q.InvertedLists[qGram]; ok {
			lists = append(lists, invertedList)
		}
	}

	for _, yPair := range MergeLists(lists) {
		yEntity := q.EntityMap[yPair.WordId]
		y := yEntity.Name

		// NOTE: special case, if delta == 0, x must be prefix of y
		if delta == 0 && !strings.HasPrefix(y, x) {
			continue
		}

		if yPair.Count < len(x)-1-delta*q.Q {
			continue
		}

		numPEDComputations += 1
		if ped := PrefixEditDistance(Normalize(x), Normalize(y), delta); ped <= delta {
			matches = append(matches, EntityPEDPair{
				Entity: yEntity,
				PED:    ped,
			})
		}
	}

	return
}

func RankMatches(matches []EntityPEDPair) (sorted []EntityPEDPair) {
	sorted = make([]EntityPEDPair, len(matches))
	copy(sorted, matches)
	sort.Slice(sorted, func(i, j int) bool {
		mi, mj := sorted[i], sorted[j]
		if mi.PED != mj.PED {
			return mi.PED < mj.PED
		}

		return mi.Entity.Score > mj.Entity.Score
	})
	return
}

var nonWordCharacterRegexp = regexp.MustCompile(`\W`)

/**
 * Normalize the given string (remove non-word characters and lower case). In
 * the lecture, this was part of the qGrams method, but we also need it as a
 * separate method when computing the EDs for the remaining candidates.
 */
func Normalize(raw string) string {
	return nonWordCharacterRegexp.ReplaceAllString(strings.ToLower(raw), "")
}

type WordIdCountPair struct {
	WordId int
	Count  int
}

// Merge the given inverted lists. Pay attention to either keep duplicates in
// the result list or keep a count of the number of each id.
//
// NOTE: It is ok, if you do this merging by simply concatenating the lists
// and then sort the concatenation. That is, you do not have to make use of
// the fact, that the lists are already sorted.
func MergeLists(lists [][]int) (ret []WordIdCountPair) {
	var idxList = make([]int, len(lists))

	checkInBound := func(i int) bool {
		return idxList[i] < len(lists[i])
	}

	checkAllInBound := func() bool {
		for i, id := range idxList {
			if id < len(lists[i]) {
				return true
			}
		}
		return false
	}

	for checkAllInBound() {
		var minWordId, count = math.MaxInt64, 0
		for i := 0; i < len(lists); i++ {
			if !checkInBound(i) {
				continue
			}
			wordId := lists[i][idxList[i]]
			if wordId < minWordId {
				minWordId = wordId
			}
		}

		for i := 0; i < len(lists); i++ {
			for checkInBound(i) && lists[i][idxList[i]] == minWordId {
				idxList[i] += 1
				count += 1
			}
		}

		ret = append(ret, WordIdCountPair{minWordId, count})
	}

	return
}

// Compute the prefix edit distance of the two given strings x and y and
// return it if it is smaller or equal to the given δ. Otherwise return δ + 1.
//
// NOTE: The method must run in time O(|x| * (|x| + δ)), as explained in the
// lecture.
//noinspection GoNilness
func PrefixEditDistance(x, y string, delta int) (ped int) {
	// NOTE: ped = 0 for empty word
	if len(x) == 0 {
		return
	}

	numRows, numCols := len(x)+1, len(x)+delta+1
	if numCols > len(y)+1 {
		numCols = len(y) + 1
	}

	var dp [][]int
	for i := 0; i < numRows; i++ {
		var row []int
		for j := 0; j < numCols; j++ {
			if j == 0 {
				row = append(row, i)
			} else if i == 0 {
				row = append(row, j)
			} else {
				row = append(row, 0)
			}
		}
		dp = append(dp, row)
	}

	ped = math.MaxInt64
	for i := 1; i < numRows; i++ {
		for j := 1; j < numCols; j++ {
			if x[i-1] == y[j-1] {
				dp[i][j] = dp[i-1][j-1]
				continue
			}

			var min = math.MaxInt64
			if dp[i-1][j-1] < min {
				min = dp[i-1][j-1]
			}
			if dp[i-1][j] < min {
				min = dp[i-1][j]
			}
			if dp[i][j-1] < min {
				min = dp[i][j-1]
			}

			dp[i][j] = min + 1
		}

		if i == numRows-1 {
			for j := 0; j < numCols; j++ {
				if dp[i][j] < ped {
					ped = dp[i][j]
				}
			}
		}
	}

	if ped > delta {
		ped = delta + 1
	}
	return
}
