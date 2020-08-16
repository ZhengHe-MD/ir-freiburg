package main

import "container/heap"

type Item struct {
	ID       int
	Priority float64
}

type ItemHeap []*Item

func (m ItemHeap) Len() int           { return len(m) }
func (m ItemHeap) Less(i, j int) bool { return m[i].Priority < m[j].Priority }
func (m ItemHeap) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }

func (m *ItemHeap) Push(x interface{}) {
	*m = append(*m, x.(*Item))
}

func (m *ItemHeap) Pop() interface{} {
	old := *m
	n := len(old)
	x := old[n-1]
	*m = old[0 : n-1]
	return x
}

type CappedItemHeap struct {
	fh  ItemHeap
	cap int
}

func NewCappedItemHeap(cap int) *CappedItemHeap {
	return &CappedItemHeap{cap: cap}
}

func (m *CappedItemHeap) Init() {
	heap.Init(&m.fh)
}

func (m *CappedItemHeap) Push(x *Item) {
	if len(m.fh) < m.cap {
		heap.Push(&m.fh, x)
		return
	}

	minItem := heap.Pop(&m.fh).(*Item)

	if minItem.Priority < x.Priority {
		heap.Push(&m.fh, x)
	} else {
		heap.Push(&m.fh, minItem)
	}
}

func (m *CappedItemHeap) Pop() (val *Item, ok bool) {
	if ok = m.fh.Len() > 0; ok {
		val = heap.Pop(&m.fh).(*Item)
	}
	return
}
