package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	inputFile, err := os.Open("./movies-benchmark.txt")
	if err != nil {
		os.Exit(-1)
		return
	}
	defer inputFile.Close()

	outputFile, err := os.Create("./movies-benchmark-minus-1.txt")
	if err != nil {
		os.Exit(-1)
		return
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	scanner.Buffer([]byte{}, 100*1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		println(line)

		parts := strings.Split(line, "\t")
		if len(parts) != 2 {
			log.Printf("invalid line %s\n", line)
			continue
		}

		query := parts[0]

		var docId int64
		var docIdStrList []string
		for _, docIdStr := range strings.Fields(parts[1]) {
			docId, err = strconv.ParseInt(docIdStr, 10, 64)
			if err != nil {
				os.Exit(-1)
				return
			}

			docIdStrList = append(docIdStrList, fmt.Sprintf("%d", docId-1))
		}

		outLine := fmt.Sprintf("%s\t%s\n", query, strings.Join(docIdStrList, " "))
		outputFile.WriteString(outLine)
	}

	fmt.Println("done")
}
