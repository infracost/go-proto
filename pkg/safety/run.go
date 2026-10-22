package safety

// RunE calls fn and returns its error, or a *Panic if it panicked. cleanup,
// if not nil, runs afterwards with the recovered panic, or nil.
func RunE(fn func() error, cleanup func(*Panic)) error {
	return wrapE(fn, cleanup)
}

// Run calls fn and returns a *Panic if it panicked, otherwise nil. cleanup,
// if not nil, runs afterwards with the recovered panic, or nil.
func Run(fn func(), cleanup func(*Panic)) error {
	return wrapE(func() error { fn(); return nil }, cleanup)
}
