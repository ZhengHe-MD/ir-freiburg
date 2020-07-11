package intersection

import . "github.com/ZhengHe-MD/ir-freiburg.git/lecture-03/postinglist"

func IntersectWithGallopingSearch(l1, l2 *PostingList) (ret *PostingList) {
	ret = NewPostingList()

	if l1.Size() > l2.Size() {
		l1, l2 = l2, l1
	}

	ret.Reserve(l1.Size())

	var i1, i2 int
	var found bool
	var endPos int
	for i1 < l1.Size() && i2 < l2.Size() {
		startPos := i2
		found, endPos = gallopingSearch(l1.GetId(i1), startPos, l2)
		if !found {
			return
		}
		found, i2 = binarySearch(l1.GetId(i1), startPos, endPos, l2)
		if found {
			ret.AddPosting(l1.GetId(i1), l1.GetScore(i1)+l2.GetScore(i2-1))
		}
		i1++
	}
	return
}

func gallopingSearch(nextId int64, startPos int, postingList *PostingList) (found bool, endPos int) {
	if postingList.GetId(startPos) > nextId {
		return true, startPos
	}

	gap := 1
	pos := startPos + gap
	for postingList.GetId(pos) < nextId {
		gap = gap + gap
		pos = pos + gap
		if pos >= postingList.Size() {
			return false, -1
		}
	}
	return true, pos
}
