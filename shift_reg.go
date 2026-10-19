package pqgram

import "iter"

type shiftReg[L Integer] struct {
	arr        []L
	head, tail uint64
}

func newShiftReg[L Integer](size uint64) *shiftReg[L] {
	return &shiftReg[L]{
		arr:  make([]L, size),
		head: 0,
		tail: size - 1,
	}
}

func (q *shiftReg[L]) Size() uint64 {
	return uint64(len(q.arr))
}

func (q *shiftReg[L]) nextIdx(idx uint64) uint64 {
	if idx >= q.Size()-1 {
		return 0
	}
	return idx + 1
}

func (q *shiftReg[L]) Shift(newVal L) L {
	head := q.arr[q.head]
	q.head = q.nextIdx(q.head)
	q.tail = q.nextIdx(q.tail)
	q.arr[q.tail] = newVal
	return head
}

func (q *shiftReg[L]) Values() iter.Seq[L] {
	return func(yield func(L) bool) {
		for _, a := range q.arr {
			if !yield(a) {
				return
			}
		}
	}
}
