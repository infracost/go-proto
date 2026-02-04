package diagnostic

import (
	"errors"

	parserpb "github.com/infracost/proto/gen/go/infracost/parser"
)

const maxNonCriticalDiagnostics = 10

// Diagnostics is a collection of Diagnostic items with filtering and aggregation methods.
type Diagnostics struct {
	list       []*Diagnostic
	critCount  int
	otherCount int
}

// NewDiagnostics creates a Diagnostics collection from a slice of diagnostics.
func NewDiagnostics(diags []*Diagnostic) *Diagnostics {
	var d Diagnostics
	for _, diag := range diags {
		d.Add(diag)
	}
	return &d
}

// FromProto creates a Diagnostics collection from protocol buffer diagnostics.
func FromProto(diags []*parserpb.Diagnostic) *Diagnostics {
	var d Diagnostics
	for _, diag := range diags {
		d.Add((*Diagnostic)(diag))
	}
	return &d
}

var ErrCriticalDiagnostics = errors.New("critical diagnostics found")

func (d *Diagnostics) String() string {
	if d == nil {
		return ""
	}
	var out string
	for _, diag := range d.list {
		out += diag.String() + "\n"
	}
	return out
}

func (d *Diagnostics) OfType(t parserpb.DiagnosticType) *Diagnostics {
	if d == nil {
		return nil
	}
	var out *Diagnostics
	for _, d := range d.list {
		if d.Type == t {
			out = out.Add(d)
		}
	}
	return out
}

func (d *Diagnostics) Merge(other *Diagnostics) *Diagnostics {
	if d == nil {
		d = &Diagnostics{}
	}
	if other == nil {
		return d
	}
	for _, diag := range other.list {
		d.Add(diag)
	}
	return d
}

func (d *Diagnostics) Unwrap() []*Diagnostic {
	if d == nil {
		return nil
	}
	return d.list
}

func (d *Diagnostics) ToProto() []*parserpb.Diagnostic {
	if d == nil {
		return nil
	}
	out := make([]*parserpb.Diagnostic, len(d.list))
	for i, diag := range d.list {
		out[i] = diag.ToProto()
	}
	return out
}

func (d *Diagnostics) Add(diags ...*Diagnostic) *Diagnostics {
	if d == nil {
		d = &Diagnostics{}
	}
	for _, diag := range diags {
		if diag.Critical {
			d.critCount++
			d.list = append(d.list, diag)
		} else if d.otherCount < maxNonCriticalDiagnostics {
			d.otherCount++
			d.list = append(d.list, diag)
		}
	}
	return d
}

func (d *Diagnostics) Len() int {
	if d == nil {
		return 0
	}
	return len(d.list)
}

func (d *Diagnostics) Critical() *Diagnostics {
	if d == nil {
		return nil
	}
	var out *Diagnostics
	for _, d := range d.list {
		if d.Critical {
			out = out.Add(d)
		}
	}
	return out
}

func (d *Diagnostics) NonCritical() *Diagnostics {
	if d == nil {
		return nil
	}
	var out *Diagnostics
	for _, d := range d.list {
		if !d.Critical && !d.Ignored {
			out = out.Add(d)
		}
	}
	return out
}

func (d *Diagnostics) Ignored() *Diagnostics {
	if d == nil {
		return nil
	}
	var out *Diagnostics
	for _, d := range d.list {
		if d.Ignored {
			out = out.Add(d)
		}
	}
	return out
}

func (d *Diagnostics) Transient() *Diagnostics {
	if d == nil {
		return nil
	}
	var out *Diagnostics
	for _, d := range d.list {
		if isTransient(d.Type) {
			out = out.Add(d)
		}
	}
	return out
}
