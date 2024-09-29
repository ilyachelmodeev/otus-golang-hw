package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("Offset over size", func(t *testing.T) {
		err := Copy("testdata/input.txt", "/tmp/copy_test_out.txt", 10_000, 0)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})
	t.Run("Offset over size", func(t *testing.T) {
		err := Copy("/dev/zero", "/tmp/copy_test_out.txt", 0, 0)
		require.ErrorIs(t, err, ErrUnsupportedFile)
	})
}
