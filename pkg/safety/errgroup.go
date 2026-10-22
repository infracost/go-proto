package safety

import "sync"

// ErrGroup runs functions concurrently and reports the first error, including
// a recovered panic, from Wait.
type ErrGroup struct {
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
}

// NewErrGroup returns an empty ErrGroup.
func NewErrGroup() *ErrGroup {
	return &ErrGroup{}
}

// Go runs f in a new goroutine. If f panics, the panic becomes the group's
// error. cleanup, if not nil, runs after f with the recovered panic, or nil.
func (e *ErrGroup) Go(f func() error, cleanup func(*Panic)) {
	e.wg.Go(func() {
		if err := wrapE(f, cleanup); err != nil {
			e.errOnce.Do(func() { e.err = err })
		}
	})
}

// Wait blocks until every function started by Go has returned, then returns
// the first non-nil error.
func (e *ErrGroup) Wait() error {
	e.wg.Wait()
	return e.err
}
