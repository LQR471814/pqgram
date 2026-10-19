package pqgram

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShiftReg(t *testing.T) {
	sr := newShiftReg[uint64](3)

	require.Equal(t, sr.Shift(1), uint64(0))
	require.Equal(t, sr.Shift(2), uint64(0))
	require.Equal(t, sr.Shift(3), uint64(0))

	require.Equal(t, sr.Shift(0), uint64(1))
	require.Equal(t, sr.Shift(0), uint64(2))
	require.Equal(t, sr.Shift(0), uint64(3))
}
