package value

import (
	"reflect"
	"slices"
	"strings"

	"github.com/infracost/go-proto/pkg/flag"
	"github.com/infracost/proto/gen/go/infracost/parser"
	prototree "github.com/infracost/proto/gen/go/infracost/tree"
)

type Valuer interface {
	ToProto() *prototree.Value
}

type Value[T Primitive] struct {
	proto *prototree.Value
}

var (
	_ Valuer   = (*Value[bool])(nil)
	_ Settable = (*Value[bool])(nil)
)

// Settable extends Valuer to allow setting the proto value via reflection.
type Settable interface {
	Valuer
	SetProto(p *prototree.Value)
}

func (v *Value[T]) SetProto(p *prototree.Value) {
	v.proto = p
}

func (v *Value[T]) IsDefault() bool {
	if v.proto == nil { // vars which were never set get flagged as default
		return true
	}
	return flag.Flags(v.proto.Flags).IsDefault()
}

func (v *Value[T]) IsDefaultOrEmpty() bool {
	return v.IsDefault() || v.IsEmpty()
}

func (v *Value[T]) Field() string {
	if v.proto.SourceFieldName == nil {
		return ""
	}
	return *v.proto.SourceFieldName
}

func (v *Value[T]) Flags() flag.Flags {
	if v.proto == nil {
		return 0
	}
	return flag.Flags(v.proto.Flags)
}

func (v *Value[T]) Source() *parser.SourceRange {
	if v.proto == nil {
		return nil
	}
	return v.proto.Source
}

type Primitive interface {
	~string | ~bool | ~int64 | ~float64 | ~uint32
}

type (
	String = Value[string]
	Bool   = Value[bool]
	Int    = Value[int64]
	Double = Value[float64]
)

func New[T Primitive](val T, flags flag.Flags, fieldName string, src *parser.SourceRange) Value[T] {
	var f *string
	if fieldName != "" {
		f = &fieldName
	}
	v := Value[T]{
		proto: &prototree.Value{
			Flags:           uint64(flags),
			SourceFieldName: f,
			Source:          src,
			Value:           nil, // populated below
		},
	}
	switch a := any(val).(type) {
	case string:
		v.proto.Value = &prototree.Value_StringValue{StringValue: a}
	case bool:
		v.proto.Value = &prototree.Value_BoolValue{BoolValue: a}
	case int64:
		v.proto.Value = &prototree.Value_IntValue{IntValue: a}
	case float64:
		v.proto.Value = &prototree.Value_DoubleValue{DoubleValue: a}
	default:
		// Named types with ~string, ~bool, ~int64, ~float64 underlying type
		rv := reflect.ValueOf(val)
		switch rv.Kind() {
		case reflect.String:
			v.proto.Value = &prototree.Value_StringValue{StringValue: rv.String()}
		case reflect.Bool:
			v.proto.Value = &prototree.Value_BoolValue{BoolValue: rv.Bool()}
		case reflect.Int64:
			v.proto.Value = &prototree.Value_IntValue{IntValue: rv.Int()}
		case reflect.Float64:
			v.proto.Value = &prototree.Value_DoubleValue{DoubleValue: rv.Float()}
		case reflect.Uint32:
			// #nosec G115
			v.proto.Value = &prototree.Value_EnumValue{EnumValue: uint32(rv.Uint())}
		}
	}
	return v
}

func FromProto[T Primitive](p *prototree.Value) Value[T] {
	return Value[T]{proto: p}
}

var EmptyString = New("", 0, "", nil)

func (v Value[T]) WithValue(value T) Value[T] {
	return New(value, v.Flags(), v.Field(), v.Source())
}

func (v Value[string]) String() string {
	return v.Value()
}

func (v Value[T]) Contains(subject string) bool {
	s, ok := any(v.Value()).(string)
	if !ok {
		return false
	}
	return strings.Contains(s, subject)
}

func (v Value[T]) HasPrefix(prefix string) bool {
	s, ok := any(v.Value()).(string)
	if !ok {
		return false
	}

	return strings.HasPrefix(s, prefix)
}

func (v Value[T]) HasSuffix(suffix string) bool {
	s, ok := any(v.Value()).(string)
	if !ok {
		return false
	}

	return strings.HasSuffix(s, suffix)
}

func (v Value[T]) EqualFold(subject string) bool {
	s, ok := any(v.Value()).(string)
	if !ok {
		return false
	}

	return strings.EqualFold(s, subject)
}

func (v Value[T]) IsOneOf(options ...T) bool {
	return slices.Contains(options, v.Value())
}

func (v Value[T]) Equal(value T) bool {
	return v.Value() == value
}

func (v Value[T]) ValueOr(fallback T) T {
	if v.IsEmpty() {
		return fallback
	}
	return v.Value()
}

func (v Value[T]) IsGreaterThan(other T) bool {
	return compare(v.Value(), other) > 0
}

func (v Value[T]) IsLessThan(other T) bool {
	return compare(v.Value(), other) < 0
}

func (v Value[T]) IsGreaterThanOrEqual(other T) bool {
	return compare(v.Value(), other) >= 0
}

func (v Value[T]) IsLessThanOrEqual(other T) bool {
	return compare(v.Value(), other) <= 0
}

// compare returns -1, 0, or 1 for ordering. Returns 0 for non-orderable types (bool).
func compare[T Primitive](a, b T) int {
	ra := reflect.ValueOf(a)
	rb := reflect.ValueOf(b)
	switch ra.Kind() {
	case reflect.String:
		as, bs := ra.String(), rb.String()
		if as < bs {
			return -1
		}
		if as > bs {
			return 1
		}
		return 0
	case reflect.Int64:
		ai, bi := ra.Int(), rb.Int()
		if ai < bi {
			return -1
		}
		if ai > bi {
			return 1
		}
		return 0
	case reflect.Float64:
		af, bf := ra.Float(), rb.Float()
		if af < bf {
			return -1
		}
		if af > bf {
			return 1
		}
		return 0
	case reflect.Uint32:
		au, bu := ra.Uint(), rb.Uint()
		if au < bu {
			return -1
		}
		if au > bu {
			return 1
		}
		return 0
	}
	return 0
}

func (v Value[T]) IsEmpty() bool {
	var empty T
	return empty == v.Value()
}

func (v Value[T]) Pointer() *T {
	x := v.Value()
	return &x
}

func (v Value[T]) Value() T {
	var zero T
	if v.proto == nil {
		return zero
	}
	rv := reflect.ValueOf(&zero).Elem()
	switch rv.Kind() {
	case reflect.String:
		rv.SetString(v.proto.GetStringValue())
	case reflect.Bool:
		rv.SetBool(v.proto.GetBoolValue())
	case reflect.Int64:
		rv.SetInt(v.proto.GetIntValue())
	case reflect.Float64:
		rv.SetFloat(v.proto.GetDoubleValue())
	case reflect.Uint32:
		rv.SetUint(uint64(v.proto.GetEnumValue()))
	}
	return zero
}

func (v Value[T]) IsFalse() bool {
	return !v.IsTrue()
}

func (v Value[T]) IsTrue() bool {
	if v.proto == nil {
		return false
	}
	return v.proto.GetBoolValue()
}

func (v Value[T]) IsSynthetic() bool {
	if v.proto == nil {
		return false
	}
	return flag.Flags(v.proto.Flags).IsSynthetic()
}

func (v Value[T]) ToProto() *prototree.Value {
	if v.proto != nil {
		return v.proto
	}
	var zero T
	return New(zero, 0, "", nil).proto
}

// List wraps a proto Value containing a ValueList.
type List[T Primitive] struct {
	items           []Value[T]
	Flags           uint64
	SourceFieldName *string
	Source          *parser.SourceRange
}

func NewList[T Primitive](items []Value[T], flags flag.Flags, fieldName string, src *parser.SourceRange) *List[T] {
	var f *string
	if fieldName != "" {
		f = &fieldName
	}
	return &List[T]{
		items:           items,
		Flags:           uint64(flags),
		SourceFieldName: f,
		Source:          src,
	}
}

func (l *List[T]) Items() []Value[T] {
	return l.items
}

func (l *List[T]) Contains(v T) bool {
	for _, item := range l.items {
		if item.Value() == v {
			return true
		}
	}
	return false
}
