package convert

import (
	"reflect"

	"github.com/infracost/go-proto/pkg/tree/value"
	prototree "github.com/infracost/proto/gen/go/infracost/tree"
)

var (
	valuerType   = reflect.TypeFor[value.Valuer]()
	settableType = reflect.TypeFor[value.Settable]()
)

// StructToValueObject converts a struct's fields into a proto ValueObject.
// Field names come from the "tree" struct tag; fields without a tag are skipped.
func StructToValueObject(v any) *prototree.ValueObject {
	return structToValueObject(reflect.ValueOf(v))
}

// ValueObjectToStruct populates a struct's fields from a proto ValueObject.
// Fields are matched using the "tree" struct tag.
func ValueObjectToStruct(obj *prototree.ValueObject, v any) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	valueObjectToStruct(obj, rv)
}

func structToValueObject(v reflect.Value) *prototree.ValueObject {
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	obj := &prototree.ValueObject{
		Entries: make(map[string]*prototree.Value),
	}

	t := v.Type()
	for i := range t.NumField() {
		field := v.Field(i)
		fieldType := t.Field(i)
		name := fieldType.Tag.Get("tree")
		if name == "" || name == "-" {
			continue
		}

		if pv := fieldToProtoValue(field); pv != nil {
			obj.Entries[name] = pv
		}
	}

	return obj
}

func fieldToProtoValue(field reflect.Value) *prototree.Value {
	if field.Type().Implements(valuerType) {
		return field.Interface().(value.Valuer).ToProto()
	}

	switch field.Kind() {
	case reflect.Pointer:
		if field.IsNil() {
			return nil
		}
		return fieldToProtoValue(field.Elem())
	case reflect.Struct:
		obj := structToValueObject(field)
		return &prototree.Value{
			Value: &prototree.Value_ObjectValue{ObjectValue: obj},
		}
	case reflect.Slice:
		list := &prototree.ValueList{
			Values: make([]*prototree.Value, field.Len()),
		}
		for i := range field.Len() {
			list.Values[i] = fieldToProtoValue(field.Index(i))
		}
		return &prototree.Value{
			Value: &prototree.Value_ListValue{ListValue: list},
		}
	}

	return nil
}

func valueObjectToStruct(obj *prototree.ValueObject, v reflect.Value) {
	t := v.Type()
	for i := range t.NumField() {
		fieldType := t.Field(i)
		name := fieldType.Tag.Get("tree")
		if name == "" || name == "-" {
			continue
		}

		entry, ok := obj.Entries[name]
		if !ok {
			continue
		}

		field := v.Field(i)
		setFieldFromProtoValue(field, entry)
	}
}

func setFieldFromProtoValue(field reflect.Value, pv *prototree.Value) {
	if reflect.PointerTo(field.Type()).Implements(settableType) {
		field.Addr().Interface().(value.Settable).SetProto(pv)
		return
	}

	switch field.Kind() {
	case reflect.Pointer:
		if pv == nil {
			return
		}
		elem := reflect.New(field.Type().Elem())
		setFieldFromProtoValue(elem.Elem(), pv)
		field.Set(elem)
	case reflect.Struct:
		if ov := pv.GetObjectValue(); ov != nil {
			valueObjectToStruct(ov, field)
		}
	case reflect.Slice:
		lv := pv.GetListValue()
		if lv == nil {
			return
		}
		elemType := field.Type().Elem()
		for _, item := range lv.Values {
			elem := reflect.New(elemType).Elem()
			setFieldFromProtoValue(elem, item)
			field.Set(reflect.Append(field, elem))
		}
	}
}
