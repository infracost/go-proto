package safety

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunE(t *testing.T) {
	t.Run("returns the function's error", func(t *testing.T) {
		want := errors.New("boom")
		var cleanupCalled bool
		err := RunE(func() error { return want }, func(p *Panic) {
			require.Nil(t, p)
			cleanupCalled = true
		})
		require.ErrorIs(t, err, want)
		require.True(t, cleanupCalled)
	})

	t.Run("recovers a panic", func(t *testing.T) {
		var cleanupPanic *Panic
		err := RunE(func() error { panic("oh no") }, func(p *Panic) { cleanupPanic = p })
		var pnc *Panic
		require.ErrorAs(t, err, &pnc)
		require.Equal(t, "oh no", pnc.Value)
		require.Same(t, pnc, cleanupPanic)
		require.Contains(t, pnc.Stack, "run_test.go")
		require.NotContains(t, pnc.Stack, "runtime/debug")
	})

	t.Run("recovers a panic in cleanup", func(t *testing.T) {
		err := RunE(func() error { return nil }, func(*Panic) { panic("cleanup") })
		require.ErrorContains(t, err, "recovered from panic while running cleanup: cleanup")
	})
}

func TestRun(t *testing.T) {
	require.NoError(t, Run(func() {}, nil))
	require.ErrorContains(t, Run(func() { panic("oh no") }, nil), "recovered from panic: oh no")
}
