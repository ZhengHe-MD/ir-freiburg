package main

import (
	"bufio"
	"fmt"
	"github.com/ZhengHe-MD/ir-freiburg.git/lecture-13/ner"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: cmd <trans-probs-filename> <word-distrib-filename>")
		os.Exit(-1)
	}

	transitionProbs, err := ner.ReadTransitionProbabilitiesFromFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	wordDistrib, err := ner.ReadWordDistributionFromFile(os.Args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	m := ner.NewNER(transitionProbs, wordDistrib)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter a sentence: ")

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		sentence := strings.Fields(line)
		wordTags, err := m.POSTag(sentence)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		namedEntities, err := m.FindNamedEntities(sentence)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		var ret strings.Builder

		ret.WriteString("POSTags: ")
		for _, wt := range wordTags {
			ret.WriteString(wt.Word)
			ret.WriteString(" \\")
			ret.WriteString(wt.Tag)
			ret.WriteByte(' ')
		}

		ret.WriteByte('\n')
		ret.WriteString("NamedEntities: ")
		for _, ne := range namedEntities {
			ret.WriteString(ne)
			ret.WriteByte(' ')
		}
		ret.WriteByte('\n')
		fmt.Println(ret.String() + "Enter another sentence: ")
	}
}
