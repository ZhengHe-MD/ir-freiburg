package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	nonAlphaCharRegex = regexp.MustCompile("[^a-zA-Z]+")
)

type InvertedIndex struct {
	invertedLists map[string][]int
}

func NewInvertedIndex() *InvertedIndex {
	return &InvertedIndex{}
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: cmd <file>")
		os.Exit(-1)
	}

	filename := os.Args[1]

	ii := NewInvertedIndex()
	err := ii.ReIndexFromFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	for word, invertedList := range ii.invertedLists {
		fmt.Printf("%s\t%d\n", word, len(invertedList))
	}
}