package pqgram

import (
	"fmt"
	"iter"
	"slices"
)

// Interval is a generic interval
type Interval struct {
	Left, Right uint64
}

func (i Interval) StrictlyContains(other Interval) bool {
	return i.Left < other.Left && i.Right > other.Right
}

func (i Interval) Size() uint64 {
	return i.Right - i.Left
}

// Node contains the data associated with a node in the tree.
//
// Interval is used for quickly checking ancestor/descendent/leaf in the tree,
// the rules are as follows:
//
//   - Leaf nodes have interval of size 1
//   - Ancestor intervals must contain descendant intervals
//   - Sibling intervals are disjoint and are ordered corresponding to the
//     order of the siblings
type Node[L Integer] struct {
	Label    L
	Interval Interval
}

// TreeCursor is an iterator over the nodes of a tree in pre-order traversal.
//
// It is analogous to a database cursor where each row is one node and rows are
// ordered in pre-order.
type TreeCursor[L Integer] = iter.Seq[Node[L]]

// Tree represents a general tree interface that can be used for traversal.
//
// V is the vertex id type, L is the label type.
type Tree[V, L Integer] interface {
	Root() V
	Vertices() iter.Seq[V]
	// A vertex label should be a value >= 1
	Label(vertex V) L
	Children(vertex V) iter.Seq[V]
}

// OutputRows represents an interface to receive the "database operations"
// which MakeTreeRows performs to obtain a valid list of rows
type OutputRows[V, L Integer] interface {
	Append(vertexId V, node *Node[L]) error
	UpdateIntervalRight(vertexId V, right uint64) error
}

// MakeTreeRows constructs a valid list of node rows (and their intervals) in
// the output by pre-order traversal of the tree.
func MakeTreeRows[V, L Integer](tree Tree[V, L], output OutputRows[V, L]) error {
	counter := uint64(0)
	return makeTreeRowsHelper(tree, output, tree.Root(), &counter)
}

func makeTreeRowsHelper[V, L Integer](
	tree Tree[V, L],
	output OutputRows[V, L],
	rootId V,
	counter *uint64,
) (err error) {
	rootNode := Node[L]{
		Label:    tree.Label(rootId),
		Interval: Interval{Left: *counter},
	}
	*counter++

	err = output.Append(rootId, &rootNode)
	if err != nil {
		return
	}

	for child := range tree.Children(rootId) {
		err = makeTreeRowsHelper(tree, output, child, counter)
		if err != nil {
			return
		}
		*counter++
	}

	err = output.UpdateIntervalRight(rootId, *counter)
	return
}

type ConcreteEdge[V Integer] struct {
	Parent, Child V
}

type ConcreteTree[V, L Integer] struct {
	vertexCount V
	labels      []L
	children    map[V][]V
}

// - root is always 0
// - vertex ids are all values [0, vertexCount)
// - labels of an index i, will correspond to vertex i
// - edges should refer vertex ids
func NewConcreteTree[V, L Integer](
	vertexCount V,
	labels []L,
	edges iter.Seq[ConcreteEdge[V]],
) (out *ConcreteTree[V, L], err error) {
	if uint64(vertexCount) != uint64(len(labels)) {
		err = fmt.Errorf(
			"number of vertices (%d) != number of labels (%d)",
			vertexCount, len(labels),
		)
		return
	}

	out = &ConcreteTree[V, L]{
		vertexCount: vertexCount,
		labels:      labels,
		children:    make(map[V][]V),
	}

	for e := range edges {
		if e.Parent >= vertexCount {
			err = fmt.Errorf(
				"parent id %d is out of range (max id: %d)",
				e.Parent, vertexCount,
			)
			return
		}
		if e.Child >= vertexCount {
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

func (t *ConcreteTree[V, L]) Root() V {
	return 0
}

func (t *ConcreteTree[V, L]) Vertices() iter.Seq[V] {
	return func(yield func(V) bool) {
		for i := V(0); i < t.vertexCount; i++ {
			if !yield(i) {
				break
			}
		}
	}
}

func (t *ConcreteTree[V, L]) Label(vertex V) L {
	return t.labels[vertex]
}

func (t *ConcreteTree[V, L]) Children(vertex V) iter.Seq[V] {
	return slices.Values(t.children[vertex])
}

type ConcreteTreeRows[V, L Integer] struct {
	rows []V
	data map[V]*Node[L]
}

func NewConcreteTreeRows[V, L Integer]() *ConcreteTreeRows[V, L] {
	return &ConcreteTreeRows[V, L]{
		data: make(map[V]*Node[L]),
	}
}

func (r *ConcreteTreeRows[V, L]) Append(vertexId V, node *Node[L]) error {
	r.rows = append(r.rows, vertexId)
	r.data[vertexId] = node
	return nil
}

func (r *ConcreteTreeRows[V, L]) UpdateIntervalRight(vertexId V, right uint64) error {
	r.data[vertexId].Interval.Right = right
	return nil
}

func (r *ConcreteTreeRows[V, L]) Cursor() TreeCursor[L] {
	return func(yield func(Node[L]) bool) {
		for _, id := range r.rows {
			if !yield(*r.data[id]) {
				break
			}
		}
	}
}
