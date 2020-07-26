package index

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type QGramIndex struct {
	Q             int
	Padding       string
	InvertedLists map[string][]int
}

// NewQGramIndex creates an empty QGramIndex.
func NewQGramIndex(q int) *QGramIndex {
	return &QGramIndex{
		Q:             q,
		Padding:       strings.Repeat("$", q-1),
		InvertedLists: make(map[string][]int),
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
		if len(parts) != 3 {
			err = errors.New("invalid line format")
			return
		}
		wordId += 1

		var word = parts[0]
		for _, qgram := range q.ComputeQGram(word) {
			q.InvertedLists[qgram] = append(q.InvertedLists[qgram], wordId)
		}
	}

	return
}


// ComputeQGram computes q-grams for padded, normalized version of given string.
func (q *QGramIndex) ComputeQGram(word string) (qGramList []string) {
	word = fmt.Sprintf("%s%s%s", q.Padding, word, q.Padding)
	for i := 0; i < len(word) - q.Q + 1; i++ {
		qGramList = append(qGramList, word[i:i+q.Q])
	}
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