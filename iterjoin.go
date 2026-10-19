package pqgram

import "iter"

func IterJoin[T any](seqs ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, sq := range seqs {
			for val := range sq {
				if !yield(val) {
					return
				}
			}
		}
	}
}
