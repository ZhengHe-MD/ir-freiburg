package ner

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	epsilon = 0.01

	BeginningTag = "BEG"
	EndTag       = "END"
)

// A simple Named Entity Recognition engine using the Viterbi algorithm.
type NER struct {
	TransitionProbs map[string]map[string]float64
	WordDistrib     map[string]map[string]float64

	_tagIndex map[string]int
	_indexTag map[int]string
}

func NewNER(transitionProbs, wordDistrib map[string]map[string]float64) *NER {
	tagIndex := make(map[string]int)
	indexTag := make(map[int]string)

	var i int
	for tag1, tag2ToProb := range transitionProbs {
		if _, ok := tagIndex[tag1]; !ok {
			tagIndex[tag1] = i
			indexTag[i] = tag1
			i += 1
		}
		for tag2 := range tag2ToProb {
			if _, ok := tagIndex[tag2]; !ok {
				tagIndex[tag2] = i
				indexTag[i] = tag2
				i += 1
			}
		}
	}

	return &NER{
		TransitionProbs: transitionProbs,
		WordDistrib:     wordDistrib,
		_tagIndex:       tagIndex,
		_indexTag: indexTag,
	}
}

type WordTag struct {
	Word string
	Tag  string
}

type Pos struct {
	x int
	y int
}

// Computes the sequence of POS-tags for the given sentence using the Viterbi
// algorithm, as explained in the lecture.
//
// The sentence is given as an array of words, without punctuation
// (e.g., commas or full stops).
//
// The transition probabilities and the word distribution for each tag are
// given in two files on the Wiki. For example code to read these files, see
// the methods given below. Note that there are two special POS-tags for the
// beginning and the end of a sentence.
//
// Returns a list of tuples (word, POS-tag) that defines the POS-tag for each
// word in the given sentence.
func (m *NER) POSTag(sentence []string) (ret []WordTag, err error) {
	numCols := len(sentence) + 2
	numRows := len(m._tagIndex)
	// init 2D matrix for dp
	dp := make([][]float64, numRows)
	for i := 0; i < numRows; i++ {
		dp[i] = make([]float64, numCols)
		if m._tagIndex[BeginningTag] == i {
			dp[i][0] = 1.0
		} else {
			dp[i][0] = 0.0
		}
	}

	// process
	trace := make(map[Pos]Pos)
	for j := 1; j < numCols; j++ {
		var word string
		if j <= len(sentence) {
			word = sentence[j-1]
		}

		for tag, i := range m._tagIndex {
			var prob, maxSourceProb float64
			var msi, msj int
			for sourceTag, k := range m._tagIndex {
				if tag == EndTag {
					prob = m.TransitionProbs[sourceTag][tag] * dp[k][j-1]
				} else {
					prob = m.TransitionProbs[sourceTag][tag] * m.WordDistrib[tag][word] * dp[k][j-1]
				}
				if prob > maxSourceProb {
					maxSourceProb = prob
					msi, msj = k, j-1
				}
			}
			dp[i][j] = maxSourceProb
			trace[Pos{i, j}] = Pos{msi, msj}
		}
	}

	// trace back
	var maxPos Pos
	var maxProb float64
	for _, i := range m._tagIndex {
		currProb := dp[i][numCols-1]
		if currProb > maxProb {
			currProb = maxProb
			maxPos.x, maxPos.y = i, numCols-1
		}
	}

	var lastPos = maxPos
	for i := len(sentence)-1; i >= 0; i-- {
		ret = append(ret, WordTag{
			Word: sentence[i],
			Tag:  m._indexTag[trace[lastPos].x],
		})
		lastPos = trace[lastPos]
	}

	// reverse
	for i := 0; i < len(ret)/2; i++ {
		j := len(ret)-i-1
		ret[i], ret[j] = ret[j], ret[i]
	}
	return
}

// Recognizes entities in the given sentence.
//
// As explained in the lecture, you can simply POS-tag the sentence using the
// method pos_tag() above and take each maximal sequence of words tagged with
// "NNP" as a named entity.
//
// Returns a list of strings that represent the recognized entities.
func (m *NER) FindNamedEntities(sentence []string) (ret []string, err error) {
	wordTags, err := m.POSTag(sentence)
	if err != nil {
		return
	}

	var li int
	var found bool
	for i, wt := range wordTags {
		if wt.Tag == "NNP" {
			if found {
				continue
			}
			found = true
			li = i
			continue
		} else {
			if found {
				parts := make([]string, 0, i-li)
				for j := li; j < i; j++ {
					parts = append(parts, wordTags[j].Word)
				}
				ret = append(ret, strings.Join(parts, " "))
				found = false
			}
		}
	}

	if found {
		parts := make([]string, 0, len(wordTags)-li)
		for j := li; j < len(wordTags); j++ {
			parts = append(parts, wordTags[j].Word)
		}
		ret = append(ret, strings.Join(parts, " "))
		found = false
	}
	return
}

// Reads the transition probabilities from the given file.
//
// The expected format of the file is one transition probability per line, in
// the format "POS-tag<TAB>POS-tag<TAB>probability"
func ReadTransitionProbabilitiesFromFile(filename string) (probs map[string]map[string]float64, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}

	probs = make(map[string]map[string]float64)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(line)

		if len(fields) != 3 {
			log.Printf("invalid line '%s' in transition probs file %s",
				line, filename)
			continue
		}

		tag1, tag2, probStr := fields[0], fields[1], fields[2]
		prob, err := strconv.ParseFloat(probStr, 64)
		if err != nil {
			log.Printf("invalid prob %f in transition probs file in line %s",
				prob, line)
			continue
		}
		if _, ok := probs[tag1]; !ok {
			probs[tag1] = make(map[string]float64)
		}
		probs[tag1][tag2] = prob
	}
	return
}

// Reads the word distribution from the given file.
//
// The expected format of the file is one word per line, in the format
// "word<TAB>POS-tag<TAB>probaility".
func ReadWordDistributionFromFile(filename string) (wordDistrib map[string]map[string]float64, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}

	wordDistrib = make(map[string]map[string]float64)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(line)

		if len(fields) != 3 {
			log.Printf("invalid line '%s' in word wordDistrib file %s",
				line, filename)
			continue
		}

		word, tag, probStr := fields[0], fields[1], fields[2]
		prob, err := strconv.ParseFloat(probStr, 64)
		if err != nil {
			log.Printf("invalid prob %f in word wordDistrib file in line %s",
				prob, line)
			continue
		}
		if _, ok := wordDistrib[tag]; !ok {
			wordDistrib[tag] = make(map[string]float64)
		}
		wordDistrib[tag][word] = prob
	}
	return
}
