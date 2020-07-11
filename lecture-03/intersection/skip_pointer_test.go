package intersection

import (
	. "github.com/ZhengHe-MD/ir-freiburg.git/lecture-03/postinglist"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntersectWithSkipPointer(t *testing.T) {
	l1, l2, l3 := NewPostingList(), NewPostingList(), NewPostingList()
	assert.NoError(t, l1.ReadFromFileWithSkipPointer("../data/example1.txt", 1))
	assert.NoError(t, l2.ReadFromFileWithSkipPointer("../data/example2.txt", 1))
	assert.NoError(t, l3.ReadFromFileWithSkipPointer("../data/example3.txt", 1))

	ret1 := IntersectWithSkipPointer(l1, l2)
	ret2 := IntersectWithSkipPointer(l1, l3)

	assert.Equal(t, "[(2, 9), (6, 5)]", ret1.String())
	assert.Equal(t, "[]", ret2.String())
}
