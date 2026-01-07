package entservice

import (
	"context"
	"os"
	"path"

	"entgo.io/ent/entc/gen"
	"github.com/dave/jennifer/jen"
	"github.com/syralon/coconutc/pkg/annotation/entproto"
)

type GenerateOption func(generator *Generator)

func WithOutput(output string) GenerateOption {
	return func(generator *Generator) {
		generator.output = output
	}
}
func WithOverwrite(overwrite bool) GenerateOption {
	return func(generator *Generator) {
		generator.overwrite = overwrite
	}
}

type Generator struct {
	output       string
	overwrite    bool
	protoPackage string
	entPackage   string

	writers []*writer
}

func NewGenerator(opts ...GenerateOption) *Generator {
	g := &Generator{}
	for _, opt := range opts {
		opt(g)
	}
	// g.writers = make([]*writer, 0)
	return g
}

func (g *Generator) Generate(ctx context.Context, graph *gen.Graph) error {
	for _, node := range graph.Nodes {
		apiOpts, err := entproto.GetAPIOptions(node.Annotations)
		if err != nil {
			return err
		}
		opts := &BuildOptions{
			APIOptions:   &apiOpts,
			ProtoPackage: g.protoPackage,
			EntPackage:   g.entPackage,
		}
		for _, w := range g.writers {
			if err := w.build(ctx, node, opts); err != nil {
				return err
			}
		}
	}
	return nil
}

type BuildOptions struct {
	APIOptions   *entproto.APIOptions
	ProtoPackage string
	EntPackage   string
}

type Builder interface {
	Build(ctx context.Context, node *gen.Type, opts *BuildOptions) (*jen.File, error)
}

type BuildFunc func(ctx context.Context, node *gen.Type, opts *BuildOptions) (*jen.File, error)

func (b BuildFunc) Build(ctx context.Context, node *gen.Type, opts *BuildOptions) (*jen.File, error) {
	return b(ctx, node, opts)
}

type writer struct {
	filename  string
	overwrite bool
	builder   Builder
}

func (w *writer) build(ctx context.Context, node *gen.Type, opts *BuildOptions) error {
	file, err := w.builder.Build(ctx, node, opts)
	if err != nil {
		return err
	}
	if !w.overwrite {
		if _, err = os.Stat(w.filename); err == nil || !os.IsNotExist(err) {
			return nil
		}
	}
	_ = os.MkdirAll(path.Dir(w.filename), 0700)
	return file.Save(w.filename)
}
