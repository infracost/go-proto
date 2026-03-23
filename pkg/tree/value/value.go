package value

import (
	"reflect"

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

func (v Value[T]) Equals(value T) bool {
	return v.Value() == value
}

func (v Value[T]) IsEmpty() bool {
	var empty T
	return empty == v.Value()
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
