package safety

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ErrGroup_Clean(t *testing.T) {
	eg := NewErrGroup()
	var x, y bool
	eg.Go(func() error {
		x = true
		return nil
	}, nil)
	eg.Go(func() error {
		y = true
		return nil
	}, nil)
	err := eg.Wait()
	require.NoError(t, err)
	require.True(t, x, "x should be true")
	require.True(t, y, "y should be true")
}

func Test_ErrGroup_Error(t *testing.T) {
	eg := NewErrGroup()
	eg.Go(func() error {
		return nil
	}, nil)
	eg.Go(func() error {
		return fmt.Errorf("error")
	}, nil)
	eg.Go(func() error {
		return nil
	}, nil)
	err := eg.Wait()
	require.Error(t, err)
}

func Test_ErrGroup_Panic(t *testing.T) {
	eg := NewErrGroup()

	var cleanup1Panic *Panic
	var cleanup2Panic *Panic
	var cleanup3Panic *Panic
	var cleanup1Called bool
	var cleanup2Called bool
	var cleanup3Called bool

	eg.Go(func() error {
		return nil
	}, func(p *Panic) {
		cleanup1Panic = p
		cleanup1Called = true
	})
	eg.Go(func() error {
		panic("oh no")
	}, func(p *Panic) {
		cleanup2Panic = p
		cleanup2Called = true
	})
	eg.Go(func() error {
		return nil
	}, func(p *Panic) {
		cleanup3Panic = p
		cleanup3Called = true
	})
	err := eg.Wait()
	require.Error(t, err)
	if !strings.HasPrefix(err.Error(), "recovered from panic: oh no\nStack Trace:\n") {
		t.Errorf("expected error to start with \"recovered from panic: oh no\nStack Trace:\n\", got %q", err.Error())
	}
	if !strings.Contains(err.Error(), "panic(") {
		t.Errorf("expected error to contain \"panic(\", got %q", err.Error())
	}
	if strings.Contains(err.Error(), pkgPath+".wrapE.func1()") {
		t.Errorf("expected error not to contain the recovering frame, got %q", err.Error())
	}
	if !strings.Contains(err.Error(), pkgPath+".wrapE(") {
		t.Errorf("expected error to keep the outer %s.wrapE frame, got %q", pkgPath, err.Error())
	}
	if strings.Contains(err.Error(), "runtime/debug") {
		t.Errorf("expected error not to contain \"runtime/debug\", got %q", err.Error())
	}
	if !strings.Contains(err.Error(), "errgroup_test.go") {
		t.Errorf("expected error to contain \"errgroup_test.go\", got %q", err.Error())
	}

	if cleanup1Panic != nil {
		t.Errorf("expected cleanup1Panic to be nil, got %v", cleanup1Panic)
	}
	if cleanup2Panic == nil {
		t.Errorf("expected cleanup2Panic to be set, got %v", cleanup2Panic)
	}
	if cleanup3Panic != nil {
		t.Errorf("expected cleanup3Panic to be nil, got %v", cleanup3Panic)
	}
	if !cleanup1Called {
		t.Errorf("expected cleanup1Called to be true, got false")
	}
	if !cleanup2Called {
		t.Errorf("expected cleanup2Called to be true, got false")
	}
	if !cleanup3Called {
		t.Errorf("expected cleanup3Called to be true, got false")
	}
}
