package intersect

import (
	. "github.com/ZhengHe-MD/ir-freiburg.git/lecture-03/postinglist"
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestIntersectBasic(t *testing.T) {
	l1, l2, l3 := NewPostingList(), NewPostingList(), NewPostingList()
	assert.NoError(t, l1.ReadFromFile("../data/example1.txt"))
	assert.NoError(t, l2.ReadFromFile("../data/example2.txt"))
	assert.NoError(t, l3.ReadFromFile("../data/example3.txt"))

	ret1 := IntersectBasic(l1, l2)
	ret2 := IntersectBasic(l1, l3)

	assert.Equal(t, "[(2, 9), (6, 5)]", ret1.String())
	assert.Equal(t, "[]", ret2.String())
}

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