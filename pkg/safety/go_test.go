package safety

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestGoRoutineSafety(t *testing.T) {

	var loggedPanic any
	var loggedStackTrace string

	waitChan := make(chan struct{})

	Go(func() {
		panic("oh no")
	}, nil, func(pnc *Panic) {
		loggedPanic = pnc.Value
		loggedStackTrace = pnc.Stack
		close(waitChan)
	})

	select {
	case <-waitChan:
		break
	case <-time.After(3 * time.Second):
		t.Fatal("timeout")
	}

	if loggedPanic != "oh no" {
		t.Errorf("expected loggedPanic to be \"oh no\", got %q", loggedPanic)
	}

	if !strings.Contains(loggedStackTrace, "panic(") {
		t.Errorf("expected loggedStackTrace to contain \"panic(\", got %q", loggedStackTrace)
	}

	if !strings.Contains(loggedStackTrace, "/go_test.go") {
		t.Errorf("expected loggedStackTrace to contain \"/go_test.go\", got %q", loggedStackTrace)
	}

	if strings.Contains(loggedStackTrace, pkgPath+".wrap.func1()") {
		t.Errorf("expected loggedStackTrace not to contain the recovering frame, got %q", loggedStackTrace)
	}

	if !strings.Contains(loggedStackTrace, pkgPath+".wrap(") {
		t.Errorf("expected loggedStackTrace to keep the outer %s.wrap frame, got %q", pkgPath, loggedStackTrace)
	}
}

func TestGoroutinePanicStillCleansUp(t *testing.T) {
	var cleanupPanic *Panic
	var loggedPanic *Panic

	waitChan := make(chan struct{})

	Go(func() {
		panic("oh no")
	}, func(p *Panic) {
		cleanupPanic = p
		close(waitChan)
	}, func(pnc *Panic) {
		loggedPanic = pnc
	})

	select {
	case <-waitChan:
		break
	case <-time.After(3 * time.Second):
		t.Fatal("timeout")
	}

	if cleanupPanic == nil {
		t.Fatal("expected cleanupPanic to be set")
	}

	if loggedPanic == nil {
		t.Fatal("expected loggedPanic to be set")
	}

	if cleanupPanic.Value != loggedPanic.Value {
		t.Errorf("expected cleanupPanic.Value to be %q, got %q", loggedPanic.Value, cleanupPanic.Value)
	}

}

func TestGoroutineNonPanicStillCleansUp(t *testing.T) {
	var cleanupPanic *Panic
	var loggedPanic *Panic
	var handlerCalled bool

	waitChan := make(chan struct{})

	Go(func() {
		fmt.Println("Hello world!")
	}, func(p *Panic) {
		cleanupPanic = p
		close(waitChan)
	}, func(pnc *Panic) {
		loggedPanic = pnc
		handlerCalled = true
	})

	select {
	case <-waitChan:
		break
	case <-time.After(3 * time.Second):
		t.Fatal("timeout")
	}

	if handlerCalled {
		t.Fatal("expected panic handler to not be called")
	}

	if cleanupPanic != nil {
		t.Errorf("expected cleanupPanic to be nil, got %v", cleanupPanic)
	}

	if loggedPanic != nil {
		t.Errorf("expected loggedPanic to be nil, got %v", loggedPanic)
	}

}

// pkgPath is spelled out rather than derived so that the trimming assertions
// fail if wrapFrame is ever derived wrongly.
const pkgPath = "github.com/infracost/go-proto/pkg/safety"

func TestWrapFrame(t *testing.T) {
	if wrapFrame != pkgPath+".wrap" {
		t.Errorf("expected wrapFrame to be %q, got %q", pkgPath+".wrap", wrapFrame)
	}
}

func TestGoNilPanicHandler(t *testing.T) {
	done := make(chan struct{})
	Go(func() {
		panic("oh no")
	}, func(*Panic) { close(done) }, nil)

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("timeout")
	}
}
