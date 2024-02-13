package queue

import (
	"log/slog"
	"sync/atomic"
	"unsafe"
)

type Node struct {
	value *SendInfo
	next  *Node
}

type LockFreeQueue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

func NewLockFreeQueue() *LockFreeQueue {
	return &LockFreeQueue{
		head: unsafe.Pointer(nil),
		tail: unsafe.Pointer(nil),
	}
}

func (l *LockFreeQueue) Enqueue(value *SendInfo) {
	node := &Node{value: value}
	for {
		if (*Node)(atomic.LoadPointer(&l.tail)) == nil {
			if atomic.CompareAndSwapPointer(&l.tail, l.tail, unsafe.Pointer(node)) {
				break
			}
		}
		if (*Node)(atomic.LoadPointer(&l.head)) == nil {
			first := (*Node)(atomic.LoadPointer(&l.tail))
			first.next = node
			if atomic.CompareAndSwapPointer(&l.head, l.head, unsafe.Pointer(node)) {
				break
			}
		}
		oldNode := (*Node)(atomic.LoadPointer(&l.head))
		oldNode.next = node
		ptr := unsafe.Pointer(node)
		ptrOldNode := unsafe.Pointer(oldNode)
		if atomic.CompareAndSwapPointer(&l.head, ptrOldNode, ptr) {
			slog.Info("Add", "node", node.value.Expression, "t", l.tail, "h", l.head)
			break
		}
	}
}

func (l *LockFreeQueue) Dequeue() (*SendInfo, bool) {
	var result *Node
	for {
		result = (*Node)(atomic.LoadPointer(&l.tail))
		if result == nil {
			return nil, false
		}
		slog.Info("Result", "res", result)
		newTail := result.next
		newTailUnsafe := unsafe.Pointer(newTail)
		if atomic.CompareAndSwapPointer(&l.tail, l.tail, newTailUnsafe) {
			break
		}
	}
	return result.value, true
}
