// Package safety runs functions and goroutines so that a panic is recovered
// and returned or handed to a callback as a *Panic, with a stack trace
// trimmed of this package's own recovery frames.
package safety
