package condqueue

import "sync"

// CondQueue is a producer wait free but consumer wait queue
// use cond instead of channel for resizeable queue.
type CondQueue[T any] struct {
	cond    *sync.Cond
	closed  bool
	buffers *List[T]
}

func NewCondQueue[T any]() *CondQueue[T] {
	return &CondQueue[T]{
		buffers: New[T](),
		cond:    sync.NewCond(&sync.Mutex{}),
	}
}

func (n *CondQueue[T]) Push(value T) {
	n.cond.L.Lock()
	if n.closed {
		n.cond.L.Unlock()
		return
	}
	n.buffers.PushBack(value)
	n.cond.L.Unlock()
	n.cond.Signal()
}

func (n *CondQueue[T]) Pop(noblock ...bool) (value T, ok bool) {
	unblock := len(noblock) > 0 && noblock[0]

	n.cond.L.Lock()
	defer n.cond.L.Unlock()
	for n.buffers.Len() == 0 {
		if n.closed || unblock {
			return
		}
		n.cond.Wait()
	}
	value = n.buffers.Remove(n.buffers.Front())
	ok = true
	return
}

func (n *CondQueue[T]) Close() {
	n.cond.L.Lock()
	n.closed = true
	n.cond.L.Unlock()
	n.cond.Broadcast()
}

func (n *CondQueue[T]) ForEach(fn func(T) bool) {
	n.cond.L.Lock()
	defer n.cond.L.Unlock()
	for e := n.buffers.Front(); e != nil; e = e.Next() {
		if fn(e.Value) {
			return
		}
	}
}

func (n *CondQueue[T]) IsClose() bool {
	n.cond.L.Lock()
	closed := n.closed
	n.cond.L.Unlock()
	return closed
}
