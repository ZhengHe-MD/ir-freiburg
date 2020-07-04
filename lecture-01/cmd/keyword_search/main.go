package main

import (
	"fmt"
	"github.com/ZhengHe-MD/ir-freiburg.git/lecture-01/index"
	"os"
)

// exercise
func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: cmd <file> <query>")
		os.Exit(-1)
	}

	filename, query := os.Args[1], os.Args[2]

	ii := index.NewInvertedIndex()
	err := ii.ReadFromFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	docIDList, err := ii.ProcessQuery(query)
	if err != nil {
		fmt.Println(err)
		return
	}

	docs := make([]string, len(docIDList))
	for i, docID := range docIDList {
		docs[i] = ii.GetDocByID(docID)
	}

	// number of docs to return
	k := 3
	for i, doc := range docs {
		if i >= k {
			break
		}
		fmt.Printf("%d %s\n", i+1, doc)
	}
	return
}
