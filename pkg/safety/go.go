package safety

// Go runs f in a new goroutine. If f panics, panicHandler receives the
// recovered panic. cleanup, if not nil, then runs in its own goroutine with
// the panic, or nil when f returned normally. Either callback may be nil.
func Go(f func(), cleanup func(*Panic), panicHandler func(pnc *Panic)) {
	go wrap(f, cleanup, panicHandler)
}
