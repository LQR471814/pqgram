package pqgram

import "iter"

type shiftReg struct {
	arr        []uint64
	head, tail uint64
}

func newShiftReg(size uint64) *shiftReg {
	return &shiftReg{
		arr:  make([]uint64, size),
		head: 0,
		tail: size - 1,
	}
}

func (q *shiftReg) Size() uint64 {
	return uint64(len(q.arr))
}

func (q *shiftReg) nextIdx(idx uint64) uint64 {
	if idx >= q.Size()-1 {
		return 0
	}
	return idx + 1
}

func (q *shiftReg) Shift(newVal uint64) uint64 {
	head := q.arr[q.head]
	q.head = q.nextIdx(q.head)
	q.tail = q.nextIdx(q.tail)
	q.arr[q.tail] = newVal
	return head
}

func (q *shiftReg) Values() iter.Seq[uint64] {
	return func(yield func(uint64) bool) {
		for _, a := range q.arr {
			if !yield(a) {
				return
			}
		}
	}
}
