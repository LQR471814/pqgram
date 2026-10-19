package pqgram

import (
	"encoding/binary"
	"fmt"
	"iter"

	"github.com/zeebo/xxh3"
)

const uint64_byte_size uint64 = 8

type labelsMultiset[L Integer] struct {
	values     map[uint64]uint64
	labelsBuff []byte
}

func newLabelsMultiset[L Integer](labelsBuffsize uint64) *labelsMultiset[L] {
	return &labelsMultiset[L]{
		values:     make(map[uint64]uint64),
		labelsBuff: make([]byte, labelsBuffsize*uint64_byte_size),
	}
}

func (m *labelsMultiset[L]) Add(labels iter.Seq[L]) {
	i := uint64(0)
	for l := range labels {
		if i >= uint64(len(m.labelsBuff)) {
			panic(fmt.Errorf("exceeded labels buffersize (%d)", len(m.labelsBuff)))
		}
		binary.BigEndian.PutUint64(m.labelsBuff[i*uint64_byte_size:], uint64(l))
		i++
	}
	key := xxh3.Hash(m.labelsBuff)
	m.values[key]++
}

func (m *labelsMultiset[L]) IntersectCardinality(other *labelsMultiset[L]) uint64 {
	card := uint64(0)

	for otherKey, otherVal := range other.values {
		val, intersecting := m.values[otherKey]
		if !intersecting {
			continue
		}
		card += min(otherVal, val)
	}

	return card
}

func (m *labelsMultiset[L]) UnionCardinality(other *labelsMultiset[L]) uint64 {
	card := uint64(0)

	// the cardinality of the union of two multisets is equal to their
	// individual cardinalities added together

	for _, v := range m.values {
		card += v
	}
	for _, v := range other.values {
		card += v
	}

	return card
}
