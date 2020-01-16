package cache

import (
	"container/heap"
	"time"
)

type timeout struct {
	expireAt time.Time
	key      string
	index    int
}

type timeoutHeap struct {
	timeouts []timeout
	indices  map[string]int
}

func newTimeoutHeap() *timeoutHeap {
	h := new(timeoutHeap)
	h.timeouts = make([]timeout, 0)
	h.indices = make(map[string]int)
	heap.Init(h)

	return h
}

func (h timeoutHeap) Len() int {
	return len(h.timeouts)
}

func (h timeoutHeap) Less(i, j int) bool {
	return h.timeouts[i].expireAt.Before(h.timeouts[j].expireAt)
}

func (h timeoutHeap) Swap(i, j int) {
	if i < 0 || j < 0 {
		return
	}

	// swap elements along with updating indicies
	h.timeouts[i], h.timeouts[j] = h.timeouts[j], h.timeouts[i]
	h.timeouts[i].index = i
	h.timeouts[j].index = j

	// update key-index mapping
	h.indices[h.timeouts[i].key] = i
	h.indices[h.timeouts[j].key] = j
}

func (h *timeoutHeap) Push(x interface{}) {
	it := x.(timeout)
	// lookup if the key is already present
	i, ok := h.indices[it.key]
	if ok {
		// if such key found
		// update the element and fix heap
		it.index = i
		h.timeouts[i] = it
		heap.Fix(h, i)
		return
	}

	// insert new element if no such key found
	n := len(h.timeouts)
	it.index = n
	h.timeouts = append(h.timeouts, it)
	h.indices[it.key] = n
}

func (h *timeoutHeap) Pop() interface{} {
	n := len(h.timeouts)
	if n == 0 {
		return nil
	}

	x := h.timeouts[n-1]
	h.timeouts = h.timeouts[0 : n-1]
	delete(h.indices, x.key)

	return x
}

func (h *timeoutHeap) take() interface{} {
	n := len(h.timeouts)
	if n == 0 {
		return nil
	}

	return h.timeouts[n-1]
}
