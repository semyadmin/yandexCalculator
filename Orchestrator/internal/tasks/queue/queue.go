package queue

import (
	"sync/atomic"
	"unsafe"
)

type Node struct {
	value any
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

// Добавляем в атомарную очередь
func (l *LockFreeQueue) Enqueue(value any) {
	node := &Node{value: value}
	for {
		tailPointer := unsafe.Pointer((*Node)(atomic.LoadPointer(&l.tail)))
		if tailPointer == nil {
			if atomic.CompareAndSwapPointer(&l.tail, tailPointer, unsafe.Pointer(node)) {
				break
			}
		}
		oldNode := (*Node)(atomic.LoadPointer(&l.head))
		ptrOldNode := unsafe.Pointer(oldNode)
		if ptrOldNode == nil {
			first := (*Node)(atomic.LoadPointer(&l.tail))
			first.next = node
			if atomic.CompareAndSwapPointer(&l.head, ptrOldNode, unsafe.Pointer(node)) {
				break
			}
		}
		oldNode.next = node
		ptr := unsafe.Pointer(node)
		if atomic.CompareAndSwapPointer(&l.head, ptrOldNode, ptr) {
			break
		}
	}
}

// Извлекаем из атомарной очереди
func (l *LockFreeQueue) Dequeue() (any, bool) {
	var result *Node
	for {
		result = (*Node)(atomic.LoadPointer(&l.tail))
		if result == nil {
			return nil, false
		}
		newTail := result.next
		newTailUnsafe := unsafe.Pointer(newTail)
		oldNode := unsafe.Pointer(result)
		if atomic.CompareAndSwapPointer(&l.tail, oldNode, newTailUnsafe) {
			if l.tail == l.head {
				l.head = unsafe.Pointer(nil)
			}
			break
		}
	}
	return result.value, true
}
