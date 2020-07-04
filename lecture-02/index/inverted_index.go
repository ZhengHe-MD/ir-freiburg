package index

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	nonAlphaCharRegex = regexp.MustCompile("[^a-zA-Z]+")
)

type Posting struct {
	DocID int64
	BM25  float64
}

type Doc struct {
	Raw string
	// document length: number of words in this doc
	DL int
}

// thread safety is not guaranteed
type InvertedIndex struct {
	invertedLists map[string][]Posting
	docs          map[int64]Doc
}

func NewInvertedIndex() *InvertedIndex {
	return &InvertedIndex{
		invertedLists: make(map[string][]Posting),
		docs:          make(map[int64]Doc),
	}
}

func (ii *InvertedIndex) GetInvertedLists() map[string][]Posting {
	return ii.invertedLists
}

func (ii *InvertedIndex) GetDocByID(id int64) Doc {
	return ii.docs[id]
}

// Construct the inverted index from the given file. The expected format of
// the file is one document per line, in the format <title>TAB<description>.
// Each entry in the inverted list associated to a word should contain a
// document id and a BM25 score. Compute the BM25 scores as follows:
//
// (1) In a first pass, compute the inverted lists with tf scores (that
//     is the number of occurrences of the word within the <title> and the
//     <description> of a document). Further, compute the document length
//     (DL) for each document (that is the number of words in the <title> and
//     the <description> of a document). Afterwards, compute the average
//     document length (AVDL).
// (2) In a second pass, iterate each inverted list and replace the tf scores
//     by BM25 scores, defined as:
//     BM25 = tf * (k+1) / (k * (1 - b + b * DL / AVDL) + tf) * log2(N/df),
//     where N is the total number of documents and df is the number of
//     documents that contains the word.
//
// On reading the file, use UTF-8 as the standard encoding. To split the
// texts into words, use the method introduced in the lecture. Make sure that
// you ignore empty words.
func (ii *InvertedIndex) ReadFromFile(filename string, bm25B, bm25K float64) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	invertedList := make(map[string][]Posting)
	docs := make(map[int64]Doc)

	scanner := bufio.NewScanner(f)
	docID, docLenSum := int64(1), 0
	for scanner.Scan() {
		// NOTE: it's counter-intuitive that the movies-benchmark counts docID from 2, not 0 or 1.
		docID += 1

		line := scanner.Text()

		wordCount := make(map[string]float64)
		words := nonAlphaCharRegex.Split(line, -1)

		docLen := 0
		for _, word := range words {
			if len(word) == 0 {
				continue
			}
			docLen += 1
			word = strings.ToLower(word)
			wordCount[word] += 1
		}
		docLenSum += docLen

		doc := Doc{
			Raw: line,
			DL:  docLen,
		}
		docs[docID] = doc

		for word, count := range wordCount {
			if _, ok := invertedList[word]; !ok {
				invertedList[word] = nil
			}
			invertedList[word] = append(invertedList[word], Posting{
				DocID: docID,
				// NOTE: tf score
				BM25: count,
			})
		}
	}

	docNum := float64(docID)
	avdl := float64(docLenSum) / docNum
	for _, postings := range invertedList {
		for i, posting := range postings {
			// BM25 = tf * (k+1) / (k * (1 - b + b * DL / AVDL) + tf) * log2(N/df)
			tf := posting.BM25

			var dl = float64(docs[posting.DocID].DL)
			idf := math.Log2(docNum / float64(len(postings)))
			if math.IsInf(bm25K, 1) {
				postings[i].BM25 = tf / (1 - bm25B + bm25B*dl/avdl) * idf
			} else {
				postings[i].BM25 = tf * (bm25K + 1) / (bm25K*(1-bm25B+bm25B*dl/avdl) + tf) * idf
			}
		}
	}

	ii.invertedLists = invertedList
	ii.docs = docs
	return
}

// getRoundedInvertedIndex round the BM25 score to 3 digits precision
func (ii *InvertedIndex) getRoundedInvertedIndex() (ret map[string][]Posting) {
	ret = make(map[string][]Posting)
	for word, Postings := range ii.invertedLists {
		for _, posting := range Postings {
			// TODO: error handling
			rounded, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", posting.BM25), 64)
			ret[word] = append(ret[word], Posting{
				DocID: posting.DocID,
				BM25:  rounded,
			})
		}
	}
	return
}

func (ii *InvertedIndex) ProcessQuery(query string) (docPostings []Posting) {
	words := nonAlphaCharRegex.Split(query, -1)
	for _, word := range words {
		if len(docPostings) == 0 {
			docPostings = ii.invertedLists[word]
		} else {
			docPostings = Merge(docPostings, ii.invertedLists[word])
		}
	}

	sort.Slice(docPostings, func(i, j int) bool {
		pi, pj := docPostings[i], docPostings[j]
		return pi.BM25 >= pj.BM25
	})

	return
}

// Compute the union of the two given inverted lists in linear time (linear
// in the total number of entries in the two lists), where the entries in
// the inverted lists are postings of form (doc_id, bm25_score) and are
// expected to be sorted by doc_id, in ascending order.
func Merge(postingsA, postingsB []Posting) (postings []Posting) {
	la, lb := len(postingsA), len(postingsB)
	i, j := 0, 0
	for i < la && j < lb {
		pi, pj := postingsA[i], postingsB[j]
		if pi.DocID == pj.DocID {
			postings = append(postings, Posting{
				DocID: pi.DocID,
				BM25:  pi.BM25 + pj.BM25,
			})
			i, j = i+1, j+1
		} else if pi.DocID <= pj.DocID {
			postings = append(postings, pi)
			i += 1
		} else {
			postings = append(postings, pj)
			j += 1
		}
	}

	if i < la {
		postings = append(postings, postingsA[i])
		i += 1
	}

	if j < lb {
		postings = append(postings, postingsB[j])
		j += 1
	}
	return
}
