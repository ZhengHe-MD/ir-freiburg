package intersection

import (
	. "github.com/ZhengHe-MD/ir-freiburg.git/lecture-03/postinglist"
	"github.com/stretchr/testify/assert"
	"testing"
)

func prepareData() (postingLists []*PostingList, err error) {
	l1, l2, l3 := NewPostingList(), NewPostingList(), NewPostingList()

	if err = l1.ReadFromFile("../data/bowling.txt"); err != nil {
		return
	}

	if err = l2.ReadFromFile("../data/film.txt"); err != nil {
		return
	}

	if err = l3.ReadFromFile("../data/rug.txt"); err != nil {
		return
	}

	postingLists = append(postingLists, l1, l2, l3)
	return
}

func prepareDataWithSentinel() (postingLists []*PostingList, err error) {
	l1, l2, l3 := NewPostingList(), NewPostingList(), NewPostingList()

	if err = l1.ReadFromFileWithSentinel("../data/bowling.txt"); err != nil {
		return
	}

	if err = l2.ReadFromFileWithSentinel("../data/film.txt"); err != nil {
		return
	}

	if err = l3.ReadFromFileWithSentinel("../data/rug.txt"); err != nil {
		return
	}

	postingLists = append(postingLists, l1, l2, l3)
	return
}

func prepareDataWithSkipPointer() (postingLists []*PostingList, err error) {
	l1, l2, l3 := NewPostingList(), NewPostingList(), NewPostingList()

	if err = l1.ReadFromFileWithSkipPointer("../data/bowling.txt", 2000); err != nil {
		return
	}

	if err = l2.ReadFromFileWithSkipPointer("../data/film.txt", 2000); err != nil {
		return
	}

	if err = l3.ReadFromFileWithSkipPointer("../data/rug.txt", 2000); err != nil {
		return
	}

	postingLists = append(postingLists, l1, l2, l3)
	return
}

func BenchmarkIntersectBasic(b *testing.B) {
	postingLists, err := prepareData()
	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		for j := 0; j < len(postingLists); j++ {
			for k := 0; k < j; k++ {
				IntersectBasic(postingLists[j], postingLists[k])
			}
		}
	}
}

func BenchmarkIntersectWithLessConditionalParts(b *testing.B) {
	postingLists, err := prepareData()
	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		for j := 0; j < len(postingLists); j++ {
			for k := 0; k < j; k++ {
				IntersectWithLessConditionalParts(postingLists[j], postingLists[k])
			}
		}
	}
}

func BenchmarkIntersectWithSentinels(b *testing.B) {
	postingLists, err := prepareDataWithSentinel()
	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		for j := 0; j < len(postingLists); j++ {
			for k := 0; k < j; k++ {
				IntersectWithSentinels(postingLists[j], postingLists[k])
			}
		}
	}
}

func BenchmarkIntersectWithBinarySearchInLongerRemainder(b *testing.B) {
	postingLists, err := prepareData()
	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		for j := 0; j < len(postingLists); j++ {
			for k := 0; k < j; k++ {
				IntersectWithBinarySearchInLongerRemainder(postingLists[j], postingLists[k])
			}
		}
	}
}

func BenchmarkIntersectWithGallopingSearch(b *testing.B) {
	postingLists, err := prepareData()
	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		for j := 0; j < len(postingLists); j++ {
			for k := 0; k < j; k++ {
				IntersectWithGallopingSearch(postingLists[j], postingLists[k])
			}
		}
	}
}

func BenchmarkIntersectWithSkipPointer(b *testing.B) {
	postingLists, err := prepareDataWithSkipPointer()
	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		for j := 0; j < len(postingLists); j++ {
			for k := 0; k < j; k++ {
				IntersectWithSkipPointer(postingLists[j], postingLists[k])
			}
		}
	}
}

func BenchmarkIntersectHybrid(b *testing.B) {
	postingLists, err := prepareData()
	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		for j := 0; j < len(postingLists); j++ {
			for k := 0; k < j; k++ {
				IntersectHybrid(postingLists[j], postingLists[k])
			}
		}
	}
}
