package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("return code", func(t *testing.T) {
		code := RunCmd([]string{"bash", "-c", "exit 5"}, nil)
		require.Equal(t, 5, code)
	})
}
