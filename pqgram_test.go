package pqgram

import (
	"math"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewProfile(t *testing.T) {
	T1, err := NewConcreteTree[uint64, uint64](
		6,
		[]uint64{'a', 'a', 'e', 'b', 'b', 'c'},
		slices.Values([]ConcreteEdge[uint64]{
			{
				Parent: 0,
				Child:  1,
			},
			{
				Parent: 0,
				Child:  4,
			},
			{
				Parent: 0,
				Child:  5,
			},
			{
				Parent: 1,
				Child:  2,
			},
			{
				Parent: 1,
				Child:  3,
			},
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	T2, err := NewConcreteTree[uint64, uint64](
		6,
		[]uint64{'a', 'a', 'e', 'b', 'b', 'x'},
		slices.Values([]ConcreteEdge[uint64]{
			{
				Parent: 0,
				Child:  1,
			},
			{
				Parent: 0,
				Child:  4,
			},
			{
				Parent: 0,
				Child:  5,
			},
			{
				Parent: 1,
				Child:  2,
			},
			{
				Parent: 1,
				Child:  3,
			},
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	t1rows := NewConcreteTreeRows[uint64, uint64]()
	err = MakeTreeRows(T1, t1rows)
	if err != nil {
		t.Fatal(err)
	}

	t2rows := NewConcreteTreeRows[uint64, uint64]()
	err = MakeTreeRows(T2, t2rows)
	if err != nil {
		t.Fatal(err)
	}

	p := uint64(2)
	q := uint64(3)

	prof1, err := NewProfile(t1rows.Cursor(), p, q)
	if err != nil {
		t.Fatal(err)
	}

	prof2, err := NewProfile(t2rows.Cursor(), p, q)
	if err != nil {
		t.Fatal(err)
	}

	dist, err := Distance(prof1, prof2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(dist)
	require.Equal(t, 0.31, math.Round(dist*100)/100, "Wrong (rounded) distance value.")
}
