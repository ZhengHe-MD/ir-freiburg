package intersection

import (
	. "github.com/ZhengHe-MD/ir-freiburg.git/lecture-03/postinglist"
	"math"
)

func IntersectWithSentinels(l1, l2 *PostingList) (ret *PostingList) {
	ret = NewPostingList()
	minSize := l1.Size()
	if l2.Size() < minSize {
		minSize = l2.Size()
	}
	ret.Reserve(minSize)

	var i1, i2 int
	for i1 < l1.Size() && i2 < l2.Size() {
		for l1.GetId(i1) < l2.GetId(i2) {
			i1++
		}

		for l2.GetId(i2) < l1.GetId(i1) {
			i2++
		}

		if l1.GetId(i1) == l2.GetId(i2) {
			if l1.GetId(i1) == math.MaxInt64 {
				break
			}
			ret.AddPosting(l1.GetId(i1), l1.GetScore(i1)+l2.GetScore(i2))
			i1++
			i2++
		}
	}
	return
}