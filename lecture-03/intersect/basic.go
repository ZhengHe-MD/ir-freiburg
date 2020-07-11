package intersect

import (
	. "github.com/ZhengHe-MD/ir-freiburg.git/lecture-03/postinglist"
)

func IntersectBasic(l1, l2 *PostingList) (ret *PostingList) {
	ret = NewPostingList()
	minSize := l1.Size()
	if l2.Size() < minSize {
		minSize = l2.Size()
	}
	ret.Reserve(minSize)

	var i1, i2 int
	for i1 < l1.Size() && i2 < l2.Size() {
		id1, id2 := l1.GetId(i1), l2.GetId(i2)
		if id1 < id2 {
			i1 += 1
		} else if id2 < id1 {
			i2 += 1
		} else {
			ret.AddPosting(id1, l1.GetScore(i1)+l2.GetScore(i2))
			i1 += 1
			i2 += 1
		}
	}
	return
}
