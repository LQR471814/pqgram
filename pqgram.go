package pqgram

import (
	"fmt"
	"iter"
)

// nullLabel returns the sentinel "null" label for a given label type.
//
// It is the zero value of the label type (0), so user-supplied labels should
// be values >= 1.
func nullLabel[L Integer]() L {
	var z L
	return z
}

type Profile[L Integer] struct {
	P, Q     uint64
	multiset *labelsMultiset[L]
}

func NewProfile[L Integer](cur TreeCursor[L], p, q uint64) (prof *Profile[L], err error) {
	prof = &Profile[L]{
		P:        p,
		Q:        q,
		multiset: newLabelsMultiset[L](p + q),
	}

	next, stop := iter.Pull(cur)
	defer stop()

	anc := newShiftReg[L](p)
	root, ok := next()
	if !ok {
		err = fmt.Errorf("no root, cursor stopped immediately")
		return
	}
	cursorNode := root
	profileHelper(&cursorNode, next, root, prof, anc)

	return
}

func profileHelper[L Integer](
	cursorNode *Node[L],
	next func() (Node[L], bool),
	root Node[L],
	prof *Profile[L],
	anc *shiftReg[L],
) bool {
	anc.Shift(root.Label)
	sib := newShiftReg[L](prof.Q)

	// if !ok, we are at the last leaf node, so cur.Interval.Size() must be 1
	cur, ok := next()
	if ok {
		*cursorNode = cur
	}

	if root.Interval.Size() == 1 {
		// root is a leaf, cur points to next node (or nil)
		prof.multiset.Add(IterJoin(anc.Values(), sib.Values()))
		return ok
	}

	// is descendent
	for root.Interval.StrictlyContains(cur.Interval) {
		sib.Shift(cur.Label)
		prof.multiset.Add(IterJoin(anc.Values(), sib.Values()))
		ok = profileHelper(cursorNode, next, cur, prof, anc)
		if !ok {
			break
		}
		cur = *cursorNode
	}

	for k := uint64(1); k <= prof.Q-1; k++ {
		sib.Shift(nullLabel[L]())
		prof.multiset.Add(IterJoin(anc.Values(), sib.Values()))
	}

	return true
}

func Distance[L Integer](a, b *Profile[L]) (out float64, err error) {
	if a.P != b.P || a.Q != b.Q {
		err = fmt.Errorf(
			"PQ parameters are mismatched (%d, %d) != (%d, %d)",
			a.P, a.Q, b.P, b.Q,
		)
		return
	}
	intersection := a.multiset.IntersectCardinality(b.multiset)
	union := a.multiset.UnionCardinality(b.multiset)
	out = 1 - 2*float64(intersection)/float64(union)
	return
}
