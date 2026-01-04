package entproto

import (
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"github.com/jhump/protoreflect/v2/protobuilder"
	cocofield "github.com/syralon/coconut/proto/syralon/coconut/field"
	"github.com/syralon/coconutc/pkg/annotation/entproto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TypeMapping interface {
	Mapping(field.Type) *protobuilder.FieldType
}

var (
	EntityTypeMapping typeMapping = map[field.Type]*protobuilder.FieldType{
		field.TypeBool:    protobuilder.FieldTypeBool(),
		field.TypeTime:    protobuilder.FieldTypeImportedMessage((&timestamppb.Timestamp{}).ProtoReflect().Descriptor()),
		field.TypeJSON:    protobuilder.FieldTypeBytes(),
		field.TypeUUID:    protobuilder.FieldTypeString(),
		field.TypeBytes:   protobuilder.FieldTypeBytes(),
		field.TypeString:  protobuilder.FieldTypeString(),
		field.TypeInt8:    protobuilder.FieldTypeInt32(),
		field.TypeInt16:   protobuilder.FieldTypeInt32(),
		field.TypeInt32:   protobuilder.FieldTypeInt32(),
		field.TypeInt:     protobuilder.FieldTypeInt64(),
		field.TypeInt64:   protobuilder.FieldTypeInt64(),
		field.TypeUint8:   protobuilder.FieldTypeUint32(),
		field.TypeUint16:  protobuilder.FieldTypeUint32(),
		field.TypeUint32:  protobuilder.FieldTypeUint32(),
		field.TypeUint:    protobuilder.FieldTypeUint64(),
		field.TypeUint64:  protobuilder.FieldTypeUint64(),
		field.TypeFloat32: protobuilder.FieldTypeFloat(),
		field.TypeFloat64: protobuilder.FieldTypeDouble(),
		// field.TypeEnum:    nil, // TODO
		// field.TypeOther:   nil, // TODO
	}

	OperationTypeMapping typeMapping = map[field.Type]*protobuilder.FieldType{
		field.TypeBool:    protobuilder.FieldTypeImportedMessage((&cocofield.BoolField{}).ProtoReflect().Descriptor()),
		field.TypeTime:    protobuilder.FieldTypeImportedMessage((&cocofield.TimestampField{}).ProtoReflect().Descriptor()),
		field.TypeJSON:    protobuilder.FieldTypeImportedMessage((&cocofield.BytesField{}).ProtoReflect().Descriptor()),
		field.TypeUUID:    protobuilder.FieldTypeImportedMessage((&cocofield.StringField{}).ProtoReflect().Descriptor()),
		field.TypeBytes:   protobuilder.FieldTypeImportedMessage((&cocofield.BytesField{}).ProtoReflect().Descriptor()),
		field.TypeEnum:    protobuilder.FieldTypeImportedMessage((&cocofield.Int32Field{}).ProtoReflect().Descriptor()),
		field.TypeString:  protobuilder.FieldTypeImportedMessage((&cocofield.StringField{}).ProtoReflect().Descriptor()),
		field.TypeInt8:    protobuilder.FieldTypeImportedMessage((&cocofield.Int32Field{}).ProtoReflect().Descriptor()),
		field.TypeInt16:   protobuilder.FieldTypeImportedMessage((&cocofield.Int32Field{}).ProtoReflect().Descriptor()),
		field.TypeInt32:   protobuilder.FieldTypeImportedMessage((&cocofield.Int32Field{}).ProtoReflect().Descriptor()),
		field.TypeInt:     protobuilder.FieldTypeImportedMessage((&cocofield.Int64Field{}).ProtoReflect().Descriptor()),
		field.TypeInt64:   protobuilder.FieldTypeImportedMessage((&cocofield.Int64Field{}).ProtoReflect().Descriptor()),
		field.TypeUint8:   protobuilder.FieldTypeImportedMessage((&cocofield.Uint32Field{}).ProtoReflect().Descriptor()),
		field.TypeUint16:  protobuilder.FieldTypeImportedMessage((&cocofield.Uint32Field{}).ProtoReflect().Descriptor()),
		field.TypeUint32:  protobuilder.FieldTypeImportedMessage((&cocofield.Uint32Field{}).ProtoReflect().Descriptor()),
		field.TypeUint:    protobuilder.FieldTypeImportedMessage((&cocofield.Uint64Field{}).ProtoReflect().Descriptor()),
		field.TypeUint64:  protobuilder.FieldTypeImportedMessage((&cocofield.Uint64Field{}).ProtoReflect().Descriptor()),
		field.TypeFloat32: protobuilder.FieldTypeImportedMessage((&cocofield.FloatField{}).ProtoReflect().Descriptor()),
		field.TypeFloat64: protobuilder.FieldTypeImportedMessage((&cocofield.DoubleField{}).ProtoReflect().Descriptor()),
		// field.TypeOther:   nil,
	}

	TypeClassicalPaginator = protobuilder.FieldTypeImportedMessage((&cocofield.ClassicalPaginator{}).ProtoReflect().Descriptor())
	TypeInfinitePaginator  = protobuilder.FieldTypeImportedMessage((&cocofield.InfinitePaginator{}).ProtoReflect().Descriptor())
)

type typeMapping map[field.Type]*protobuilder.FieldType

func (m typeMapping) Mapping(t field.Type) *protobuilder.FieldType {
	return m[t]
}

func NewField(name string, field *gen.Field, mapping TypeMapping) (*protobuilder.FieldBuilder, error) {
	fieldOpts, err := entproto.GetFieldOptions(field.Annotations)
	if err != nil {
		return nil, err
	}
	t := fieldOpts.Type
	if t == 0 {
		t = field.Type.Type
	}
	fi := protobuilder.NewField(protoreflect.Name(name), mapping.Mapping(t))
	if fieldOpts.TypeRepeated {
		fi.SetRepeated()
	}
	return fi, nil
}

func MustNewField(name string, field *gen.Field, mapping TypeMapping) *protobuilder.FieldBuilder {
	fi, err := NewField(name, field, mapping)
	if err != nil {
		panic(err)
	}
	return fi
}

func PaginatorType(style entproto.PaginatorStyle) *protobuilder.FieldType {
	if style == entproto.InfinitePaginator {
		return TypeInfinitePaginator
	}
	return TypeClassicalPaginator
}
