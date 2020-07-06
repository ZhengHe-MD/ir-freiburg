package main

import (
	"fmt"
	"github.com/ZhengHe-MD/ir-freiburg.git/lecture-02/evaluator"
	"github.com/ZhengHe-MD/ir-freiburg.git/lecture-02/index"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: cmd <dataset> <benchmark>")
		os.Exit(-1)
	}

	datasetFilename, benchmarkFilename := os.Args[1], os.Args[2]
	options := index.RefinementOptions{ExcludingStopWords: true}

	ii := index.NewInvertedIndex()
	//err := ii.ReadFromFile(datasetFilename, 0.75, 1.25)
	//err := ii.ReadFromFile(datasetFilename, 0.1, 0.75)
	//err := ii.ReadFromFile(datasetFilename, 0.3, 1.0)
	//err := ii.ReadFromFile(datasetFilename, 0.34, 1.35)
	err := ii.ReadFromFile(datasetFilename, 0.11, 0.77, options)
	if err != nil {
		log.Println(err)
		return
	}

	benchmark, err := evaluator.ReadBenchmark(benchmarkFilename)
	if err != nil {
		log.Println(err)
		return
	}

	mps := evaluator.Evaluate(ii, benchmark, options)
	fmt.Printf("MP@3: %.3f\n", mps.MPAt3)
	fmt.Printf("MP@R: %.3f\n", mps.MPAtR)
	fmt.Printf("MAP: %.3f\n", mps.MAP)
}
