package main

import (
	"fmt"
	"github.com/ZhengHe-MD/ir-freiburg.git/lecture-01/index"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: cmd <file>")
		os.Exit(-1)
	}

	filename := os.Args[1]

	ii := index.NewInvertedIndex()
	err := ii.ReadFromFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	for word, invertedList := range ii.GetInvertedLists() {
		fmt.Printf("%s\t%d\n", word, len(invertedList))
	}
}
