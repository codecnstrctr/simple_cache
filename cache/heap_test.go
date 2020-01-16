package cache

import (
	"container/heap"
	"testing"
	"time"
)

func TestTimeoutHeap(t *testing.T) {
	h := newTimeoutHeap()

	item := heap.Pop(h)
	if item != nil {
		t.Error("err on popping from empty heap")
	}

	foo := timeout{
		expireAt: time.Now().Add(time.Minute * 3),
		key:      "foo",
	}

	bar := timeout{
		expireAt: time.Now().Add(time.Minute),
		key:      "bar",
	}

	baz := timeout{
		expireAt: time.Now().Add(time.Hour),
		key:      "baz",
	}

	heap.Push(h, foo)
	heap.Push(h, bar)
	heap.Push(h, baz)
	item = heap.Pop(h)
	if item.(timeout).key != "bar" {
		t.Errorf("err on pop: got %s <> want %s", item.(timeout).key, "bar")
	}

	baz.expireAt = time.Now()
	heap.Push(h, baz)
	item = heap.Pop(h)
	if item.(timeout).key != "baz" {
		t.Errorf("err on pop: got %s <> want %s", item.(timeout).key, "baz")
	}
}
