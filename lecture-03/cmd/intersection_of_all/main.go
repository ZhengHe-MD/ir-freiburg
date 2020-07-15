package main

import (
	"fmt"
	"github.com/ZhengHe-MD/ir-freiburg.git/lecture-03/intersection"
	. "github.com/ZhengHe-MD/ir-freiburg.git/lecture-03/postinglist"
	"os"
)

func prepareData() (postingLists []*PostingList, err error) {
	l1, l2, l3 := NewPostingList(), NewPostingList(), NewPostingList()

	if err = l1.ReadFromFile("./data/bowling.txt"); err != nil {
		return
	}

	if err = l2.ReadFromFile("./data/film.txt"); err != nil {
		return
	}

	if err = l3.ReadFromFile("./data/rug.txt"); err != nil {
		return
	}

	postingLists = append(postingLists, l1, l2, l3)
	return
}

func main() {
	postingLists, err := prepareData()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	var ret *PostingList
	for i := 0; i < len(postingLists); i++ {
		if i == 0 {
			ret = postingLists[i]
		} else {
			ret = intersection.IntersectBasic(ret, postingLists[i])
		}
	}

	ret.Iterate(func(i int) {
		fmt.Println(ret.GetId(i), "\t", ret.GetScore(i))
	})
}
