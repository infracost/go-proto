// Package log is the plugin side of the plugin stderr log contract. A plugin
// writes one "[LEVEL] message" line per entry to stderr; the host's go-plugin
// client reads each line and recovers the level from the "[LEVEL]" prefix, so
// the prefix must start the line and carry no timestamp (the host adds one).
//
// Calls below the configured level build no string and write nothing. The
// level comes from LevelEnv, then the deprecated DeprecatedLevelEnv, and
// defaults to warn.
package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

const (
	// LevelEnv names the variable the host sets to pass its log level to a
	// plugin.
	LevelEnv = "LOG_LEVEL"
	// DeprecatedLevelEnv is read when LevelEnv is unset.
	//
	// Deprecated: hosts should set LevelEnv.
	DeprecatedLevelEnv = "INFRACOST_CLI_LOG_LEVEL"
)

// Level is a log severity, ordered from most to least verbose.
type Level int32

// Levels understood by the host.
const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
)

var prefixes = [...]string{
	LevelTrace: "[TRACE] ",
	LevelDebug: "[DEBUG] ",
	LevelInfo:  "[INFO] ",
	LevelWarn:  "[WARN] ",
	LevelError: "[ERROR] ",
}

// ParseLevel maps a level name to a Level, case-insensitively. Fatal and panic
// map to LevelError; unknown names map to LevelWarn.
func ParseLevel(s string) Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "trace":
		return LevelTrace
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "error", "fatal", "panic":
		return LevelError
	default:
		return LevelWarn
	}
}

// LevelFromEnv returns the level set by the host, or LevelWarn when it set
// none.
func LevelFromEnv() Level {
	if v, ok := os.LookupEnv(LevelEnv); ok {
		return ParseLevel(v)
	}
	if v, ok := os.LookupEnv(DeprecatedLevelEnv); ok {
		return ParseLevel(v)
	}
	return LevelWarn
}

// Logger writes leveled lines in the host's format. It is safe for concurrent
// use.
type Logger struct {
	mu    sync.Mutex
	out   io.Writer
	level atomic.Int32
}

// New returns a Logger writing entries at level or above to out.
func New(out io.Writer, level Level) *Logger {
	l := &Logger{out: out}
	l.level.Store(int32(level))
	return l
}

// SetLevel changes the minimum level written.
func (l *Logger) SetLevel(level Level) { l.level.Store(int32(level)) }

// Enabled reports whether an entry at level would be written.
func (l *Logger) Enabled(level Level) bool { return level >= Level(l.level.Load()) }

func (l *Logger) logf(level Level, format string, args ...any) {
	if !l.Enabled(level) {
		return
	}
	line := prefixes[level] + fmt.Sprintf(format, args...) + "\n"
	l.mu.Lock()
	defer l.mu.Unlock()
	_, _ = io.WriteString(l.out, line)
}

// Tracef logs at LevelTrace.
func (l *Logger) Tracef(format string, args ...any) { l.logf(LevelTrace, format, args...) }

// Debugf logs at LevelDebug.
func (l *Logger) Debugf(format string, args ...any) { l.logf(LevelDebug, format, args...) }

// Infof logs at LevelInfo.
func (l *Logger) Infof(format string, args ...any) { l.logf(LevelInfo, format, args...) }

// Warnf logs at LevelWarn.
func (l *Logger) Warnf(format string, args ...any) { l.logf(LevelWarn, format, args...) }

// Errorf logs at LevelError.
func (l *Logger) Errorf(format string, args ...any) { l.logf(LevelError, format, args...) }

var std = New(os.Stderr, LevelFromEnv())

// Default returns the process-wide Logger, which writes to stderr at the level
// from the environment.
func Default() *Logger { return std }

// SetLevel changes the minimum level of the Default logger.
func SetLevel(level Level) { std.SetLevel(level) }

// Tracef logs at LevelTrace on the Default logger.
func Tracef(format string, args ...any) { std.logf(LevelTrace, format, args...) }

// Debugf logs at LevelDebug on the Default logger.
func Debugf(format string, args ...any) { std.logf(LevelDebug, format, args...) }

// Infof logs at LevelInfo on the Default logger.
func Infof(format string, args ...any) { std.logf(LevelInfo, format, args...) }

// Warnf logs at LevelWarn on the Default logger.
func Warnf(format string, args ...any) { std.logf(LevelWarn, format, args...) }

// Errorf logs at LevelError on the Default logger.
func Errorf(format string, args ...any) { std.logf(LevelError, format, args...) }
