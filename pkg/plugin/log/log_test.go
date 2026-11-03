package log

import (
	"bytes"
	"os"
	"testing"
)

func TestLoggerFormatAndFilter(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf, LevelInfo)

	l.Tracef("t")
	l.Debugf("d")
	l.Infof("i %d", 1)
	l.Warnf("w %s", "x")
	l.Errorf("e")

	want := "[INFO] i 1\n[WARN] w x\n[ERROR] e\n"
	if got := buf.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	buf.Reset()
	l.SetLevel(LevelTrace)
	l.Tracef("t")
	if got := buf.String(); got != "[TRACE] t\n" {
		t.Errorf("after SetLevel got %q", got)
	}
}

type countingStringer struct{ n *int }

func (c countingStringer) String() string { *c.n++; return "" }

func TestFilteredCallDoesNotFormat(t *testing.T) {
	var buf bytes.Buffer
	n := 0
	New(&buf, LevelWarn).Debugf("%s", countingStringer{&n})
	if n != 0 || buf.Len() != 0 {
		t.Errorf("filtered call formatted %d times and wrote %q", n, buf.String())
	}
}

func TestParseLevel(t *testing.T) {
	for in, want := range map[string]Level{
		"trace":   LevelTrace,
		" DEBUG ": LevelDebug,
		"info":    LevelInfo,
		"warn":    LevelWarn,
		"warning": LevelWarn,
		"error":   LevelError,
		"fatal":   LevelError,
		"panic":   LevelError,
		"":        LevelWarn,
		"bogus":   LevelWarn,
	} {
		if got := ParseLevel(in); got != want {
			t.Errorf("ParseLevel(%q) = %d, want %d", in, got, want)
		}
	}
}

func TestLevelFromEnv(t *testing.T) {
	unset := func(k string) {
		if v, ok := os.LookupEnv(k); ok {
			t.Cleanup(func() { _ = os.Setenv(k, v) })
			_ = os.Unsetenv(k)
		}
	}
	unset(LevelEnv)
	unset(DeprecatedLevelEnv)
	if got := LevelFromEnv(); got != LevelWarn {
		t.Errorf("no env: got %d, want warn", got)
	}

	t.Setenv(DeprecatedLevelEnv, "debug")
	if got := LevelFromEnv(); got != LevelDebug {
		t.Errorf("deprecated env: got %d, want debug", got)
	}

	t.Setenv(LevelEnv, "error")
	if got := LevelFromEnv(); got != LevelError {
		t.Errorf("%s should win: got %d, want error", LevelEnv, got)
	}
}
