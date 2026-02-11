package entproto

import (
	"errors"

	"entgo.io/ent/schema/field"
	"github.com/syralon/coconutc/pkg/annotation/helper"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

type Filter uint16

const (
	FilterEQ = 1 << iota
	FilterNE
	FilterGT
	FilterGTE
	FilterLT
	FilterLTE
	FilterBETWEEN
	FilterIN
	FilterAll = 1<<iota - 1
)

func (f Filter) Filters() []Filter {
	return bitTwiddling(f)
}

type FieldOptions struct {
	//Name           string

	Filterable     bool
	Immutable      bool
	Settable       bool
	Sensitive      bool
	Filter         Filter
	Orderable      bool
	Type           field.Type
	TypeRepeated   bool
	ProtoEnum      bool
	ProtoEnumValue map[string]int32
}

type fieldAnnotation struct {
	FieldOptions
}

func (a *fieldAnnotation) Name() string { return fieldAnnotationName }

type FieldOption func(*fieldAnnotation)

//func WithFieldName(name string) FieldOption {
//	return func(a *fieldAnnotation) {
//		a.FieldOptions.Name = name
//	}
//}

func WithFieldImmutable(immutable bool) FieldOption {
	return func(a *fieldAnnotation) {
		a.Immutable = immutable
	}
}

func WithFieldSettable(settable bool) FieldOption {
	return func(a *fieldAnnotation) {
		a.Settable = settable
	}
}

func WithFieldFilterable(filterable bool) FieldOption {
	return func(a *fieldAnnotation) {
		a.Filterable = filterable
	}
}

// WithFieldSensitive
// The sensitive field will not be shown.
func WithFieldSensitive(sensitive bool) FieldOption {
	return func(a *fieldAnnotation) {
		a.Sensitive = sensitive
	}
}

func WithFieldOrderable(orderable bool) FieldOption {
	return func(a *fieldAnnotation) {
		a.Orderable = orderable
	}
}

func WithFieldFilter(filters ...Filter) FieldOption {
	return func(a *fieldAnnotation) {
		a.Filter = 0
		for _, f := range filters {
			a.Filter |= f
		}
	}
}

func WithFieldType(fieldType field.Type, repeated ...bool) FieldOption {
	return func(a *fieldAnnotation) {
		a.Type = fieldType
		a.TypeRepeated = len(repeated) > 0 && repeated[0]
	}
}

func WithFieldProtoEnum(isEnum bool, values map[string]int32) FieldOption {
	return func(a *fieldAnnotation) {
		a.ProtoEnum = isEnum
		a.ProtoEnumValue = values
	}
}

func Field(opts ...FieldOption) entc.Annotation {
	a := &fieldAnnotation{
		FieldOptions: defaultFieldOption,
	}
	for _, option := range opts {
		option(a)
	}
	return a
}

var defaultFieldOption = FieldOptions{
	Filterable: true,
	Filter:     FilterAll,
}

func GetFieldOptions(annotations gen.Annotations) (FieldOptions, error) {
	s := &fieldAnnotation{}
	err := helper.GetAnnotations(annotations, fieldAnnotationName, s)
	if errors.Is(err, helper.ErrAnnotationNotFound) {
		return defaultFieldOption, nil
	}
	if err != nil {
		return FieldOptions{}, err
	}
	return s.FieldOptions, nil
}
