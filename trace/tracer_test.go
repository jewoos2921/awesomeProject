package trace

import (
	"bytes"
	"io"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Errorf("Return from New should not be in nil")
	} else {
		tracer.Trace("Hello trace package")
		if buf.String() != "Hello trace package\n" {
			t.Errorf("Trace should not write '%s'.", buf.String())
		}
	}
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}
