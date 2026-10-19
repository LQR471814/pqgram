package pqgram

import (
	"fmt"
	"iter"
)

type Edge struct {
	Parent, Child uint64
}

type Tree struct {
	vertexCount uint64
	labels      []uint64
	children    map[uint64][]uint64
}

// - vertex ids are all values [1, vertexCount]
// - edges should refer vertex ids
// - labels of an index i, will correspond to vertex i+1
func NewTree(
	vertexCount uint64,
	labels []uint64,
	edges iter.Seq[Edge],
) (out *Tree, err error) {
	if vertexCount != uint64(len(labels)) {
		err = fmt.Errorf(
			"number of vertices (%d) != number of labels (%d)",
			vertexCount, len(labels),
		)
		return
	}

	out = &Tree{
		vertexCount: vertexCount,
		labels:      labels,
		children:    make(map[uint64][]uint64),
	}

	for e := range edges {
		if e.Parent == 0 {
			err = fmt.Errorf("cannot use sentinel id 0 in edge parent")
			return
		}
		if e.Child == 0 {
			err = fmt.Errorf("cannot use sentinel id 0 in edge child")
			return
		}
		if e.Parent > vertexCount {
			err = fmt.Errorf(
				"parent id %d is out of range (max id: %d)",
				e.Parent, vertexCount,
			)
			return
		}
		if e.Child > vertexCount {
			err = fmt.Errorf(
				"child id %d is out of range (max id: %d)",
				e.Child, vertexCount,
			)
			return
		}
		out.children[e.Parent] = append(out.children[e.Parent], e.Child)
	}

	return
}

func (t *Tree) GetLabel(id uint64) uint64 {
	return t.labels[id-1]
}

func (t *Tree) Children(id uint64) iter.Seq[uint64] {
	return func(yield func(uint64) bool) {
		for _, child := range t.children[id] {
			if !yield(child) {
				return
			}
		}
	}
}
