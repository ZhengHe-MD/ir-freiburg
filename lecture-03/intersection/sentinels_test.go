package intersection

import (
	. "github.com/ZhengHe-MD/ir-freiburg.git/lecture-03/postinglist"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntersectWithSentinels(t *testing.T) {
	l1, l2, l3 := NewPostingList(), NewPostingList(), NewPostingList()
	assert.NoError(t, l1.ReadFromFileWithSentinel("../data/example1.txt"))
	assert.NoError(t, l2.ReadFromFileWithSentinel("../data/example2.txt"))
	assert.NoError(t, l3.ReadFromFileWithSentinel("../data/example3.txt"))

	ret1 := IntersectWithSentinels(l1, l2)
	ret2 := IntersectWithSentinels(l1, l3)

	assert.Equal(t, "[(2, 9), (6, 5)]", ret1.String())
	assert.Equal(t, "[]", ret2.String())
}
