package pqgram

import (
	"slices"
	"testing"
)

func TestNewProfile(t *testing.T) {
	T1, err := NewTree(
		6,
		[]uint64{'a', 'a', 'e', 'b', 'b', 'c'},
		slices.Values([]Edge{
			{
				Parent: 1,
				Child:  2,
			},
			{
				Parent: 1,
				Child:  5,
			},
			{
				Parent: 1,
				Child:  6,
			},
			{
				Parent: 2,
				Child:  3,
			},
			{
				Parent: 2,
				Child:  4,
			},
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	T2, err := NewTree(
		6,
		[]uint64{'a', 'a', 'e', 'b', 'b', 'x'},
		slices.Values([]Edge{
			{
				Parent: 1,
				Child:  2,
			},
			{
				Parent: 1,
				Child:  5,
			},
			{
				Parent: 1,
				Child:  6,
			},
			{
				Parent: 2,
				Child:  3,
			},
			{
				Parent: 2,
				Child:  4,
			},
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	p := uint64(2)
	q := uint64(3)
	prof1 := NewProfile(T1, p, q, 1)
	prof2 := NewProfile(T2, p, q, 1)
	dist, err := Distance(prof1, prof2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(dist)
}
