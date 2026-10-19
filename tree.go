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
type Node struct {
	Label    uint64
	Interval Interval
}

// TreeCursor is an iterator over the nodes of a tree in pre-order traversal.
//
// It is analogous to a database cursor where each row is one node and rows are
// ordered in pre-order.
type TreeCursor = iter.Seq[Node]

// Tree represents a general tree interface that can be used for traversal
type Tree interface {
	Root() uint64
	Vertices() iter.Seq[uint64]
	// A vertex label should be a value >= 1
	Label(vertex uint64) uint64
	Children(vertex uint64) iter.Seq[uint64]
}

// OutputRows represents an interface to receive the "database operations"
// which MakeTreeRows performs to obtain a valid list of rows
type OutputRows interface {
	Append(vertexId uint64, node *Node) error
	UpdateIntervalRight(vertexId, right uint64) error
}

// MakeTreeRows constructs a valid list of node rows (and their intervals) in
// the output by pre-order traversal of the tree.
func MakeTreeRows(tree Tree, output OutputRows) error {
	counter := uint64(0)
	return makeTreeRowsHelper(tree, output, tree.Root(), &counter)
}

func makeTreeRowsHelper(
	tree Tree,
	output OutputRows,
	rootId uint64,
	counter *uint64,
) (err error) {
	rootNode := Node{
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

type ConcreteEdge struct {
	Parent, Child uint64
}

type ConcreteTree struct {
	vertexCount uint64
	labels      []uint64
	children    map[uint64][]uint64
}

// - root is always 0
// - vertex ids are all values [0, vertexCount)
// - labels of an index i, will correspond to vertex i
// - edges should refer vertex ids
func NewConcreteTree(
	vertexCount uint64,
	labels []uint64,
	edges iter.Seq[ConcreteEdge],
) (out *ConcreteTree, err error) {
	if vertexCount != uint64(len(labels)) {
		err = fmt.Errorf(
			"number of vertices (%d) != number of labels (%d)",
			vertexCount, len(labels),
		)
		return
	}

	out = &ConcreteTree{
		vertexCount: vertexCount,
		labels:      labels,
		children:    make(map[uint64][]uint64),
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

func (t *ConcreteTree) Root() uint64 {
	return 0
}

func (t *ConcreteTree) Vertices() iter.Seq[uint64] {
	return func(yield func(uint64) bool) {
		for i := uint64(0); i < t.vertexCount; i++ {
			if !yield(i) {
				break
			}
		}
	}
}

func (t *ConcreteTree) Label(vertex uint64) uint64 {
	return t.labels[vertex]
}

func (t *ConcreteTree) Children(vertex uint64) iter.Seq[uint64] {
	return slices.Values(t.children[vertex])
}

type ConcreteTreeRows struct {
	rows []uint64
	data map[uint64]*Node
}

func NewConcreteTreeRows() *ConcreteTreeRows {
	return &ConcreteTreeRows{
		data: make(map[uint64]*Node),
	}
}

func (r *ConcreteTreeRows) Append(vertexId uint64, node *Node) error {
	r.rows = append(r.rows, vertexId)
	r.data[vertexId] = node
	return nil
}

func (r *ConcreteTreeRows) UpdateIntervalRight(vertexId, right uint64) error {
	r.data[vertexId].Interval.Right = right
	return nil
}

func (r *ConcreteTreeRows) Cursor() TreeCursor {
	return func(yield func(Node) bool) {
		for _, id := range r.rows {
			if !yield(*r.data[id]) {
				break
			}
		}
	}
}
