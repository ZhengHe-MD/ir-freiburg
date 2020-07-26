package main

import (
	"bufio"
	"fmt"
	"github.com/ZhengHe-MD/ir-freiburg.git/lecture-05/index"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: cmd <file>")
		os.Exit(-1)
	}

	qi := index.NewQGramIndex(3)
	err := qi.BuildFromFile(os.Args[1])
	if err != nil {
		fmt.Printf("BuildFromFile err %v", err)
		os.Exit(-1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter a query:")

	for scanner.Scan() {
		startTime := time.Now()

		line := scanner.Text()
		if line == "EOF" {
			fmt.Println("See you.")
			break
		}
		x := line
		delta := len(index.Normalize(x)) / 4
		fmt.Printf("x: %s delta: %d\n", x, delta)

		matches, numPEDComputations := qi.FindMatches(x, delta)
		sortedMatches := index.RankMatches(matches)

		for i, match := range sortedMatches {
			if i >= 5 {
				break
			}
			fmt.Printf("%s\t%s\n", match.Entity.Name, match.Entity.Description)
		}

		if len(sortedMatches) > 5 {
			fmt.Printf("total #match: %d\n", len(sortedMatches))
		}

		fmt.Printf("total #PEDComputations: %d\n", numPEDComputations)
		fmt.Printf("query time: %v\n", time.Now().Sub(startTime))
		fmt.Println("Please enter a query:")
	}

	return
}
