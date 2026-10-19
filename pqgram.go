package pqgram

import (
	"fmt"
	"iter"
)

const NULL_LABEL uint64 = 0

type Profile struct {
	P, Q     uint64
	multiset *labelsMultiset
}

func NewProfile(cur TreeCursor, p, q uint64) (prof *Profile, err error) {
	prof = &Profile{
		P:        p,
		Q:        q,
		multiset: newLabelsMultiset(p + q),
	}

	next, stop := iter.Pull(cur)
	defer stop()

	anc := newShiftReg(p)
	root, ok := next()
	if !ok {
		err = fmt.Errorf("no root, cursor stopped immediately")
		return
	}
	cursorNode := root
	profileHelper(&cursorNode, next, root, prof, anc)

	return
}

func profileHelper(
	cursorNode *Node,
	next func() (Node, bool),
	root Node,
	prof *Profile,
	anc *shiftReg,
) bool {
	anc.Shift(root.Label)
	sib := newShiftReg(prof.Q)

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
		sib.Shift(NULL_LABEL)
		prof.multiset.Add(IterJoin(anc.Values(), sib.Values()))
	}

	return true
}

func Distance(a, b *Profile) (out float64, err error) {
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
