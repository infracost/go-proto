package safety

import (
	"fmt"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
)

// Panic is a recovered panic value and the stack it was raised on.
type Panic struct {
	Value any
	Stack string
}

// Error formats the panic value and its stack trace.
func (p *Panic) Error() string {
	return fmt.Sprintf("recovered from panic: %v\nStack Trace:\n%s", p.Value, p.Stack)
}

func wrapE(f func() error, cleanup func(*Panic)) (err error) {
	if f == nil {
		return
	}
	defer func() {
		var pnc *Panic
		if r := recover(); r != nil {
			pnc = &Panic{Value: r, Stack: simplifyStackTrace()}
			err = pnc
		}
		if cleanup != nil {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("recovered from panic while running cleanup: %v\nStack Trace:\n%s\nOriginal panic: %#v", r, simplifyStackTrace(), err)
				}
			}()
			cleanup(pnc)
		}
	}()
	err = f()
	return err
}

func cleanupPanicHandler(originalHandler func(*Panic)) func(pnc *Panic) {
	return func(pnc *Panic) {
		err := fmt.Errorf("recovered from panic while running cleanup: %v\nStack Trace:\n%s", pnc.Value, pnc.Stack)
		if originalHandler != nil {
			originalHandler(&Panic{Value: err, Stack: simplifyStackTrace()})
		}
	}
}

func wrap(f func(), cleanup func(*Panic), handler func(pnc *Panic)) {
	if f == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			pnc := &Panic{Value: r, Stack: simplifyStackTrace()}
			if handler != nil {
				handler(pnc)
			}
			if cleanup != nil {
				go wrap(func() { cleanup(pnc) }, nil, cleanupPanicHandler(handler))
			}
		} else if cleanup != nil {
			go wrap(func() {
				cleanup(nil)
			}, nil, cleanupPanicHandler(handler))
		}
	}()
	f()
}

// wrapFrame is wrap's fully qualified name as it appears in a stack trace,
// for example "example.com/mod/safety.wrap".
var wrapFrame string

func init() {
	wrapFrame = runtime.FuncForPC(reflect.ValueOf(wrap).Pointer()).Name()
}

// simplifyStackTrace removes the stack trace lines from this package and runtime.debug for simpler debugging
func simplifyStackTrace() (trace string) {

	defer func() {
		if r := recover(); r != nil {
			trace = ""
		}
	}()

	stack := debug.Stack()

	if !strings.Contains(string(stack), wrapFrame) {
		return string(stack)
	}

	input := strings.Split(string(stack), "\n")
	lines := make([]string, 0, len(input))

	var hide bool
	var hideNext bool

	for _, line := range input {
		if hideNext {
			hideNext = false
			continue
		}
		if line == "runtime/debug.Stack()" {
			hide = true
			hideNext = true
			continue
		}
		if hide {
			if strings.Contains(line, wrapFrame) {
				hide = false
				hideNext = true
			}
			continue
		}
		lines = append(lines, line)
	}

	trace = strings.Join(lines, "\n")
	return trace
}
