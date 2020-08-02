package index

import (
	"bufio"
	"fmt"
	"github.com/james-bowman/sparse"
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
	Score float64
}

type Doc struct {
	Raw string
	// document length: number of terms in this doc
	DL int
}

// thread safety is not guaranteed
type InvertedIndex struct {
	invertedLists map[string][]Posting
	docs          map[int64]Doc
	numTerms      int
	numDocs       int
	terms         []string
	termToIdx     map[string]int
	tdMatrix      *sparse.DOK
}

func NewInvertedIndex() *InvertedIndex {
	return &InvertedIndex{
		invertedLists: make(map[string][]Posting),
		docs:          make(map[int64]Doc),
		termToIdx:     make(map[string]int),
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
// Each entry in the inverted list associated to a term should contain a
// document id and a Score score. Compute the Score scores as follows:
//
// (1) In a first pass, compute the inverted lists with tf scores (that
//     is the number of occurrences of the term within the <title> and the
//     <description> of a document). Further, compute the document length
//     (DL) for each document (that is the number of terms in the <title> and
//     the <description> of a document). Afterwards, compute the average
//     document length (AVDL).
// (2) In a second pass, iterate each inverted list and replace the tf scores
//     by Score scores, defined as:
//     Score = tf * (k+1) / (k * (1 - b + b * DL / AVDL) + tf) * log2(N/df),
//     where N is the total number of documents and df is the number of
//     documents that contains the term.
//
// On reading the file, use UTF-8 as the standard encoding. To split the
// texts into terms, use the method introduced in the lecture. Make sure that
// you ignore empty terms.
func (ii *InvertedIndex) ReadFromFile(filename string, bm25B, bm25K float64, options RefinementOptions) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	invertedList := make(map[string][]Posting)
	docs := make(map[int64]Doc)

	scanner := bufio.NewScanner(f)
	docID, docLenSum := int64(0), 0
	for scanner.Scan() {
		docID += 1

		line := scanner.Text()

		termCount := make(map[string]float64)
		terms := nonAlphaCharRegex.Split(line, -1)

		docLen := 0
		for _, term := range terms {
			if len(term) == 0 {
				continue
			}

			if options.ExcludingStopWords && IsStopWord(term) {
				continue
			}
			docLen += 1
			term = strings.ToLower(term)
			termCount[term] += 1
			if _, ok := invertedList[term]; !ok {
				invertedList[term] = nil
				ii.terms = append(ii.terms, term)
				ii.termToIdx[term] = len(ii.terms) - 1
			}
		}
		docLenSum += docLen

		doc := Doc{
			Raw: line,
			DL:  docLen,
		}
		docs[docID] = doc

		for term, count := range termCount {
			invertedList[term] = append(invertedList[term], Posting{
				DocID: docID,
				Score: count,
			})
		}
	}

	docNum := float64(docID)
	avdl := float64(docLenSum) / docNum
	for _, postings := range invertedList {
		for i, posting := range postings {
			// Score = tf * (k+1) / (k * (1 - b + b * DL / AVDL) + tf) * log2(N/df)
			tf := posting.Score
			idf := math.Log2(docNum / float64(len(postings)))

			switch options.RankingScore {
			case RankingScoreTF:
				continue
			case RankingScoreTFIDF:
				postings[i].Score = tf * idf
			case RankingScoreBM25:
				var dl = float64(docs[posting.DocID].DL)
				if math.IsInf(bm25K, 1) {
					postings[i].Score = tf / (1 - bm25B + bm25B*dl/avdl) * idf
				} else {
					postings[i].Score = tf * (bm25K + 1) / (bm25K*(1-bm25B+bm25B*dl/avdl) + tf) * idf
				}
			case RankingscoreBM25WithoutIDF:
				var dl = float64(docs[posting.DocID].DL)
				if math.IsInf(bm25K, 1) {
					postings[i].Score = tf / (1 - bm25B + bm25B*dl/avdl)
				} else {
					postings[i].Score = tf * (bm25K + 1) / (bm25K*(1-bm25B+bm25B*dl/avdl) + tf)
				}
			}
		}
	}

	// build term-document matrix
	ii.numDocs = int(docID)
	ii.numTerms = len(invertedList)

	ii.invertedLists = invertedList
	ii.docs = docs
	return
}

// getRoundedInvertedIndex round the Score score to 3 digits precision, only for testing purpose
func (ii *InvertedIndex) getRoundedInvertedIndex() (ret map[string][]Posting) {
	ret = make(map[string][]Posting)
	for term, Postings := range ii.invertedLists {
		for _, posting := range Postings {
			// TODO: error handling
			rounded, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", posting.Score), 64)
			ret[term] = append(ret[term], Posting{
				DocID: posting.DocID,
				Score: rounded,
			})
		}
	}
	return
}

type Normalization int

const (
	None Normalization = iota
	L1
	L2
)

// preprocessingVSM preprocess the invertedList to build a term-document matrix
func (ii *InvertedIndex) PreprocessingVSM(normalization Normalization) {
	ii.tdMatrix = sparse.NewDOK(ii.numTerms, ii.numDocs)

	switch normalization {
	case None:
		for termID, term := range ii.terms {
			for _, posting := range ii.invertedLists[term] {
				docID := int(posting.DocID - 1)
				ii.tdMatrix.Set(termID, docID, posting.Score)
			}
		}
	case L1:
		normalizer := make(map[int]float64, ii.numDocs)
		for termID, term := range ii.terms {
			for _, posting := range ii.invertedLists[term] {
				docID := int(posting.DocID - 1)
				ii.tdMatrix.Set(termID, docID, posting.Score)
				normalizer[docID] += posting.Score
			}
		}

		ii.tdMatrix.DoNonZero(func(i, j int, v float64) {
			ii.tdMatrix.Set(i, j, v/normalizer[j])
		})
	case L2:
		normalizer := make(map[int]float64, ii.numDocs)
		for termID, term := range ii.terms {
			for _, posting := range ii.invertedLists[term] {
				docID := int(posting.DocID - 1)
				ii.tdMatrix.Set(termID, docID, posting.Score)
				normalizer[docID] += posting.Score * posting.Score
			}
		}

		for docID, squareSum := range normalizer {
			normalizer[docID] = math.Sqrt(squareSum)
		}

		ii.tdMatrix.DoNonZero(func(i, j int, v float64) {
			ii.tdMatrix.Set(i, j, v/normalizer[j])
		})
	}

	return
}

// getRoundedTDMatrix round the term-document matrix element to 3 digits precision, only for testing purpose
func (ii *InvertedIndex) getRoundedTDMatrix() (matrix *sparse.DOK) {
	numRows, numCols := ii.tdMatrix.Dims()
	matrix = sparse.NewDOK(numRows, numCols)
	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			rounded, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", ii.tdMatrix.At(i, j)), 64)
			matrix.Set(i, j, rounded)
		}
	}
	return
}

func (ii *InvertedIndex) ProcessQueryVSM(query string, options RefinementOptions) (docPostings []Posting) {
	terms := nonAlphaCharRegex.Split(query, -1)

	qv := sparse.NewDOK(1, ii.numTerms)
	for _, term := range terms {
		if len(term) == 0 {
			continue
		}

		if options.ExcludingStopWords && IsStopWord(term) {
			continue
		}

		if idx, ok := ii.termToIdx[term]; ok {
			qv.Set(0, idx, qv.At(0, idx)+1)
		}
	}

	var productCSR sparse.CSR
	productCSR.Mul(qv, ii.tdMatrix)

	productCSR.DoNonZero(func(i, j int, v float64) {
		docPostings = append(docPostings, Posting{
			DocID: int64(j) + 1,
			Score: v,
		})
	})

	sort.Slice(docPostings, func(i, j int) bool {
		pi, pj := docPostings[i], docPostings[j]
		return pi.Score >= pj.Score
	})

	return
}
