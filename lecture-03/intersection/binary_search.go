package intersection

import (
	. "github.com/ZhengHe-MD/ir-freiburg.git/lecture-03/postinglist"
)

func IntersectWithBinarySearchInLongerRemainder(l1, l2 *PostingList) (ret *PostingList) {
	ret = NewPostingList()

	if l1.Size() > l2.Size() {
		l1, l2 = l2, l1
	}

	ret.Reserve(l1.Size())

	var i1, i2 int
	var found bool
	for i1 < l1.Size() && i2 < l2.Size() {
		found, i2 = binarySearch(l1.GetId(i1), i2, l2.Size()-1, l2)
		if found {
			ret.AddPosting(l1.GetId(i1), l1.GetScore(i1)+l2.GetScore(i2-1))
		}
		i1++
	}
	return
}

func binarySearch(id int64, startPos, endPos int, postingList *PostingList) (found bool, nextPos int) {
	li, ri := startPos, endPos

	if id > postingList.GetId(ri) {
		return false, ri + 1
	}

	if id < postingList.GetId(li) {
		return false, li
	}

	for li <= ri {
		mi := (li + ri) / 2
		if postingList.GetId(mi) == id {
			return true, mi + 1
		} else if postingList.GetId(mi) < id {
			li = mi + 1
		} else {
			ri = mi - 1
		}
	}

	if li > ri {
		li, ri = ri, li
	}
	return false, ri
}
