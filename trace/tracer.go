package trace

import (
	"fmt"
	"io"
)

/*
Tracer is the interface that describes an object capable of
tracing events throught code
*/
type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

//New returns a new tracer which writes to the given output
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}

//Off returns a new tracer that does nothing when calling the Trace method
func Off() Tracer {
	return &nilTracer{}
}
