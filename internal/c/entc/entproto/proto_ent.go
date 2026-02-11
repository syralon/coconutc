package entproto

import (
	"fmt"
	"path"

	"entgo.io/ent/entc/gen"
	"github.com/iancoleman/strcase"
	"github.com/jhump/protoreflect/v2/protobuilder"
	"github.com/syralon/coconutc/pkg/annotation/entproto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type EntBuildOption interface {
	applyEnt(*EntBuilder)
}

type EntBuilder struct {
	options
}

func NewEntBuilder(options ...EntBuildOption) *EntBuilder {
	eb := &EntBuilder{}
	for _, option := range options {
		option.applyEnt(eb)
	}
	return eb
}

func (b *EntBuilder) Build(ctx *Context, graph *gen.Graph) ([]*protobuilder.FileBuilder, error) {
	file := ctx.NewFile(path.Join(b.path, "ent.proto"), b.protoPackage, b.goPackage)
	var messages []*protobuilder.MessageBuilder
	for _, node := range graph.Nodes {
		enums, err := b.enums(ctx, node)
		if err != nil {
			return nil, err
		}
		for _, e := range enums {
			file.AddEnum(e)
		}
		message := ctx.NewMessage(node.Name)
		file.AddMessage(message)
		messages = append(messages, message)
	}
	h := NewMessageBuildHelper(WithSkipFunc(func(f *gen.Field, opt entproto.FieldOptions) bool { return opt.Sensitive }))
	for i, node := range graph.Nodes {
		if err := h.Build(ctx, messages[i], node); err != nil {
			return nil, err
		}
	}
	return []*protobuilder.FileBuilder{file}, nil
}

func (b *EntBuilder) enums(ctx *Context, node *gen.Type) ([]*protobuilder.EnumBuilder, error) {
	var enums []*protobuilder.EnumBuilder
	for _, fi := range node.Fields {
		opt, err := entproto.GetFieldOptions(fi.Annotations)
		if err != nil {
			return nil, err
		}
		if !opt.ProtoEnum {
			continue
		}
		en := ctx.NewEnum(fmt.Sprintf("%s_%s", strcase.ToScreamingSnake(node.Name), strcase.ToScreamingSnake(fi.Name)))
		en.AddValue(protobuilder.NewEnumValue(protoreflect.Name(string(en.Name()) + "_UNSPECIFIED")))
		for k, v := range opt.ProtoEnumValue {
			val := protobuilder.NewEnumValue(protoreflect.Name(string(en.Name()) + "_" + strcase.ToScreamingSnake(k)))
			if v > 0 {
				val.SetNumber(protoreflect.EnumNumber(v))
			}
			en.AddValue(val)
		}
		enums = append(enums, en)
	}
	return enums, nil
}
