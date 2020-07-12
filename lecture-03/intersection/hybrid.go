package intersection

import . "github.com/ZhengHe-MD/ir-freiburg.git/lecture-03/postinglist"

func IntersectHybrid(l1, l2 *PostingList) (ret *PostingList) {
	if l1.Size() < l2.Size() {
		l1, l2 = l2, l1
	}
	if l1.Size()/l2.Size() < 20 {
		return IntersectWithLessConditionalParts(l1, l2)
	} else {
		return IntersectWithGallopingSearch(l1, l2)
	}
}
