package index

// head -20 words+frequencies.txt
// the     514438
// a       323284
// and     315883
// film    291714
// by      260597
// is      246859
// of      186705
// in      150309
// directed        150099
// to      91407
// was     90158
// s       73651
// it      73139
// on      72879
// written 61775
// as      51536
// for     49760
// with    40801
// drama   38745
// short   35074

var stopWordMap = map[string]bool{
	"the":      true,
	"a":        true,
	"and":      true,
	"film":     true,
	"by":       true,
	"is":       true,
	"of":       true,
	"in":       true,
	"directed": true,
	"to":       true,
	"was":      true,
	"s":        true,
	"it":       true,
	"on":       true,
	"written":  true,
	"as":       true,
	"for":      true,
	"with":     true,
	"drama":    true,
	"short":    true,
}

type RankingScore int

const (
	RankingScoreTF             RankingScore = 1
	RankingScoreTFIDF                       = 2
	RankingScoreBM25                        = 3
	RankingScoreBM25WithoutIDF              = 4
)

type RefinementOptions struct {
	ExcludingStopWords bool
	RankingScore       RankingScore
}

func IsStopWord(word string) bool {
	return stopWordMap[word]
}
