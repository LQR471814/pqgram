package pqgram

import "fmt"

type Profile struct {
	P, Q     uint64
	multiset *labelsMultiset
}

func NewProfile(tree *Tree, p, q, root uint64) (prof *Profile) {
	prof = &Profile{
		P:        p,
		Q:        q,
		multiset: newLabelsMultiset(p + q),
	}

	anc := newShiftReg(p)
	profileHelper(tree, prof, root, anc)
	return
}

func profileHelper(tree *Tree, prof *Profile, root uint64, anc *shiftReg) {
	anc.Shift(tree.GetLabel(root))
	sib := newShiftReg(prof.Q)

	children := tree.children[root]

	if len(children) == 0 {
		// root is a leaf
		prof.multiset.Add(IterJoin(anc.Values(), sib.Values()))
		return
	}

	for _, child := range children {
		sib.Shift(tree.GetLabel(child))
		prof.multiset.Add(IterJoin(anc.Values(), sib.Values()))
		profileHelper(tree, prof, child, anc)
	}

	for k := uint64(1); k <= prof.Q-1; k++ {
		sib.Shift(0)
		prof.multiset.Add(IterJoin(anc.Values(), sib.Values()))
	}
}

func Distance(a, b *Profile) (out float64, err error) {
	if a.P != b.P || a.Q != b.Q {
		err = fmt.Errorf("PQ parameters are mismatched (%d, %d) != (%d, %d)", a.P, a.Q, b.P, b.Q)
		return
	}
	intersection := a.multiset.IntersectCardinality(b.multiset)
	union := a.multiset.UnionCardinality(b.multiset)
	out = 1 - 2*float64(intersection)/float64(union)
	return
}
