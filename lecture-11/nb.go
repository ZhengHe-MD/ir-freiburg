package main

import (
	"bufio"
	"fmt"
	"github.com/james-bowman/sparse"
	"gonum.org/v1/gonum/mat"
	"math"
	"os"
	"regexp"
	"strings"
)

var (
	nonAlphaCharRegex = regexp.MustCompile("[^a-zA-Z]+")
)

// GenerateVocabularies Reads the given file and generates vocabularies mapping from label/class
// to label ids and from word to word id.
//
// You should call this ONLY on your training data.
func GenerateVocabularies(filename string) (wordVocabulary, classVocabulary map[string]int, err error) {
	// Map from label/class to label id.
	wordVocabulary = make(map[string]int)
	// Map from word to word id.
	classVocabulary = make(map[string]int)

	classID, wordID := 0, 0

	// Read the file (containing the training data)
	f, err := os.Open(filename)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(strings.TrimSpace(line), "\t")

		if len(parts) != 2 {
			err = fmt.Errorf("invalid line %s", line)
			return
		}

		label, text := parts[0], parts[1]

		if _, ok := classVocabulary[label]; !ok {
			classVocabulary[label] = classID
			classID += 1
		}

		words := strings.Fields(nonAlphaCharRegex.ReplaceAllString(strings.ToLower(text), " "))

		for _, word := range words {
			if _, ok := wordVocabulary[word]; !ok {
				wordVocabulary[word] = wordID
				wordID += 1
			}
		}
	}

	err = scanner.Err()
	return
}

// ReadLabeledData Reads the given file and returns a sparse document-term matrix as well as a
// list of labels of each document. You need to provide a class and word vocabulary. Words not
// in the vocabulary are ignored. Documents labeled with classes not in the class vocabulary
// are also ignored.
//
// The returned document-term matrix X has size n x m, where n is the number of documents and m
// the number of word ids. The value at i, j denotes the number of times word id j is present in
// document i.
//
// The returned labels vector y has size n (one label for each document). The value at index j
// denotes the label (class id) of document j.
func ReadLabeledData(filename string, classVocabulary, wordVocabulary map[string]int) (X *sparse.DOK, y []int, err error) {

	f, err := os.Open(filename)
	if err != nil {
		return
	}

	var labels, rows, cols []int
	var values []float64
	var currRowID = 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		line := scanner.Text()
		parts := strings.Split(strings.TrimSpace(line), "\t")
		if len(parts) != 2 {
			err = fmt.Errorf("invalid line %s", line)
			return
		}

		label, text := parts[0], parts[1]
		if _, ok := classVocabulary[label]; ok {
			labels = append(labels, classVocabulary[label])
			words := strings.Fields(nonAlphaCharRegex.ReplaceAllString(strings.ToLower(text), " "))
			for _, w := range words {
				if _, ok := wordVocabulary[w]; ok {
					wordID := wordVocabulary[w]
					rows = append(rows, currRowID)
					cols = append(cols, wordID)
					values = append(values, 1.0)
				}
			}
		}

		currRowID += 1
	}

	if err = scanner.Err(); err != nil {
		return
	}

	X = sparse.NewDOK(currRowID, len(wordVocabulary))
	for i := 0; i < len(rows); i++ {
		ri, ci, vi := rows[i], cols[i], values[i]
		X.Set(ri, ci, X.At(ri, ci)+vi)
	}

	y = labels
	return
}

const epsilon = 0.1

type NaiveBayes struct {
	PC  []float64
	PWC []float64
}

// Trains on the sparse document-term matrix X and associated labels y.
//
// In the test case below, p_wc is a class-term-matrix and has a row
// for each class and a column for each term. So the value at i,j is
// the p_wc for the j-th term in the i-th class.
//
// p_c is an array of global probabilities for each class.
//
// Remember to use epsilon = 1/10 for your training, as described in the
// lecture!
func (m *NaiveBayes) Train(X *sparse.DOK, y []int, classVocab, wordVocab map[string]int) (err error) {
	tcm := make(map[int]int)
	for _, classID := range y {
		tcm[classID] += 1
	}

	// NOTE: depends on class id that should be a series of continuous natural numbers
	// starting from zero.
	pc := make([]float64, len(classVocab))
	for classID, tc := range tcm {
		pc[classID] = float64(tc) / float64(len(y))
	}

	numClass := len(classVocab)
	numTerms := len(wordVocab)
	pwc := make([]float64, numClass*numTerms)
	nc := make([]float64, numClass)

	X.DoNonZero(func(i, j int, v float64) {
		cid, tid := y[i], j
		pwc[cid*numTerms+tid] += v
		nc[cid] += v
	})

	for i := 0; i < numClass; i++ {
		for j := 0; j < numTerms; j++ {
			pwc[i*numTerms+j] = (pwc[i*numTerms+j] + epsilon) / (nc[i] + epsilon*float64(numTerms))
		}
	}

	m.PWC = pwc
	m.PC = pc
	return
}

// Predicts a label for each example in the document-term matrix,
// based on the learned probabilities stored in this class.
//
// Returns a list of predicted label ids.
func (m *NaiveBayes) Predict(Xt *sparse.DOK) (yt []int, err error) {
	var cd sparse.CSR

	pwcLog := make([]float64, len(m.PWC))
	for i, w := range m.PWC {
		if w > 0 {
			pwcLog[i] = math.Log2(w)
		}
	}
	// class-term matrix * term-document matrix => class-document matrix
	cd.Mul(mat.NewDense(len(m.PC), len(m.PWC)/len(m.PC), pwcLog), Xt.T())

	numClass, numDocument := cd.Dims()
	for i := 0; i < numDocument; i++ {
		var predictClassID int
		var maxProb = -math.MaxFloat64
		for j := 0; j < numClass; j++ {
			if cd.At(j, i) > maxProb {
				maxProb = cd.At(j, i)
				predictClassID = j
			}
		}
		yt = append(yt, predictClassID)
	}

	return
}

const evaluateEpsilon = 0.0000001

// Predicts the labels of X and computes the precisions, recalls and
// F1 scores for each class.
func (m *NaiveBayes) Evaluate(Xt *sparse.DOK, y []int) (precisionList, recallList, f1List []float64, err error) {
	yt, err := m.Predict(Xt)
	if err != nil {
		return
	}

	numClass := len(m.PC)

	var precision, recall, f1 float64
	for i := 0; i < numClass; i++ {
		var numDocLabeled, numDocClassified, numTruePositive int
		for j := 0; j < len(y); j++ {
			if y[j] == i {
				numDocLabeled += 1
			}

			if yt[j] == i {
				numDocClassified += 1
				if yt[j] == y[j] {
					numTruePositive += 1
				}
			}
		}

		precision = float64(numTruePositive) / (float64(numDocClassified) + evaluateEpsilon)
		recall = float64(numTruePositive) / (float64(numDocLabeled) + evaluateEpsilon)
		f1 = 2 * precision * recall / (precision + recall + evaluateEpsilon)
		precisionList = append(precisionList, precision)
		recallList = append(recallList, recall)
		f1List = append(f1List, f1)
	}

	return
}

func LoadStopWordSet(filename string) (stopWordSet map[string]bool, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}

	stopWordSet = make(map[string]bool)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		word := strings.TrimSpace(line)
		if word == "" {
			continue
		}
		stopWordSet[word] = true
	}

	err = scanner.Err()
	return
}

func reverseVocab(vocab map[string]int) (reversedVocab map[int]string) {
	reversedVocab = make(map[int]string)
	for k, v := range vocab {
		reversedVocab[v] = k
	}
	return
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run xx.go <train-input> <test-input>")
		os.Exit(-1)
	}

	trainInput := os.Args[1]
	testInput := os.Args[2]

	wordVocab, classVocab, err := GenerateVocabularies(trainInput)
	if err != nil {
		fmt.Printf("GenerateVocabularies err %v\n", err)
		os.Exit(-1)
	}

	nb := &NaiveBayes{}

	X, y, err := ReadLabeledData(trainInput, classVocab, wordVocab)
	if err != nil {
		fmt.Printf("ReadLabeledData err %v\n", err)
		os.Exit(-1)
	}

	Xt, yt, err := ReadLabeledData(testInput, classVocab, wordVocab)
	if err != nil {
		fmt.Printf("ReadLabeledData err %v\n", err)
		os.Exit(-1)
	}

	if err = nb.Train(X, y, classVocab, wordVocab); err != nil {
		fmt.Printf("train err %v\n", err)
		os.Exit(-1)
	}

	// Predict labels for the test dataset
	// Output the precision, recall, and F1 score on the test data, for
	// each class separately as well as the (unweighted) average over all
	// classes.
	precisionList, recallList, f1List, err := nb.Evaluate(Xt, yt)
	if err != nil {
		fmt.Printf("evaluate err %v\n", err)
		os.Exit(-1)
	}
	fmt.Printf("precision: %v\n", precisionList)
	fmt.Printf("recall: %v\n", recallList)
	fmt.Printf("f1 score: %v\n", f1List)

	var averageF1 float64
	for _, f1 := range f1List {
		averageF1 += f1
	}
	averageF1 = averageF1 / float64(len(f1List))
	fmt.Printf("average f1 score: %v\n", averageF1)

	stopWordSet, err := LoadStopWordSet("./data/stopwords.txt")
	if err != nil {
		fmt.Printf("load stop word set err %v\n", err)
		os.Exit(-1)
	}

	reversedWordVocab, reversedClassVocab := reverseVocab(wordVocab), reverseVocab(classVocab)

	// Print the 20 words with the highest p_wc values per class which do
	// not appear in the stopwords.txt provided on the Wiki.
	numTerms := len(nb.PWC) / len(nb.PC)
	for i := 0; i < len(nb.PC); i++ {
		h := NewCappedItemHeap(20)
		for j := 0; j < numTerms; j++ {
			word := reversedWordVocab[j]
			if stopWordSet[word] {
				continue
			}
			h.Push(&Item{ID: j, Priority: nb.PWC[i*numTerms+j]})
		}

		var words []string
		for {
			item, ok := h.Pop()
			if !ok {
				break
			}

			words = append(words, reversedWordVocab[item.ID])
		}

		fmt.Printf("%s\t", reversedClassVocab[i])
		for i := len(words) - 1; i >= 0; i-- {
			fmt.Printf("%s ", words[i])
		}
		fmt.Printf("\n")
	}
}
