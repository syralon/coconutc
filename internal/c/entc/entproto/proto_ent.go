package entproto

import (
	"path"

	"entgo.io/ent/entc/gen"
	"github.com/jhump/protoreflect/v2/protobuilder"
	"github.com/syralon/coconutc/pkg/annotation/entproto"
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
