package entservice

import (
	"context"
	"os"
	"path"
	"strings"

	"entgo.io/ent/entc/gen"
	"github.com/dave/jennifer/jen"
	"github.com/syralon/coconutc/internal/tools/text"
	"github.com/syralon/coconutc/pkg/annotation/entproto"
)

type GenerateOption func(generator *Generator)

func WithOverwrite(overwrite bool) GenerateOption {
	return func(generator *Generator) {
		generator.overwrite = overwrite
	}
}

func WithEntPath(entPath string) GenerateOption {
	return func(generator *Generator) {
		generator.entPath = entPath
	}
}

func WithProtoPath(protoPath string) GenerateOption {
	return func(generator *Generator) {
		generator.protoPath = protoPath
	}
}

func WithOutput(output string) GenerateOption {
	return func(generator *Generator) {
		generator.output = output
	}
}

type Generator struct {
	output    string
	entPath   string
	protoPath string
	overwrite bool

	entPackage   string
	protoPackage string
	module       string

	writers []*writer
}

func NewGenerator(opts ...GenerateOption) (*Generator, error) {
	g := &Generator{entPath: "ent"}
	for _, opt := range opts {
		opt(g)
	}

	if err := g.init(); err != nil {
		return nil, err
	}
	return g, nil
}

func (g *Generator) init() (err error) {
	if g.output == "" {
		g.output = "."
	}
	g.module, err = text.Module(g.output)
	if err != nil {
		return err
	}
	if g.protoPath == "" {
		g.protoPath = text.ProtoModule(g.module)
	}
	g.protoPackage = path.Join(g.module, g.protoPath)
	g.entPackage = path.Join(g.module, g.entPath)

	var entityPath = "internal/domain/entity"
	var repositoryPath = "internal/domain/repository"
	var dataPath = "internal/infra/data"
	var servicePath = "internal/transport/service"
	var txPath = "internal/infra/tx"
	g.addWriter(
		func(n *gen.Type) string { return path.Join(g.output, entityPath, strings.ToLower(n.Name)) + ".go" },
		EntityBuilder("entity", true),
	).addWriter(
		func(n *gen.Type) string { return path.Join(g.output, repositoryPath, strings.ToLower(n.Name)) + ".go" },
		RepositoryInterfaceBuilder("repository", path.Join(g.module, entityPath)),
	).addWriter(
		func(n *gen.Type) string { return path.Join(g.output, dataPath, strings.ToLower(n.Name)) + ".go" },
		RepositoryBuilder("data", path.Join(g.module, entityPath), path.Join(g.module, txPath)),
	).addWriter(
		func(n *gen.Type) string {
			return path.Join(g.output, servicePath, strings.ToLower(n.Name)+"service", strings.ToLower(n.Name)) + ".go"
		},
		ServiceBuilder(path.Join(g.module, entityPath)),
	)
	return nil
}

func (g *Generator) addWriter(filename func(n *gen.Type) string, builder Builder) *Generator {
	g.writers = append(g.writers, &writer{
		filename:  filename,
		builder:   builder,
		overwrite: g.overwrite,
	})
	return g
}

func (g *Generator) Generate(ctx context.Context, graph *gen.Graph) error {
	data := &RenderData{
		Module:       g.module,
		ProtoPath:    g.protoPath,
		ProtoPackage: g.protoPackage,
		Services:     make([]string, 0, len(graph.Nodes)),
		overwrite:    g.overwrite,
	}

	for _, node := range graph.Nodes {
		data.Services = append(data.Services, node.Name)

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
			if err = w.build(ctx, node, opts); err != nil {
				return err
			}
		}
	}
	return data.RenderAllFile(g.output)
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
	filename  func(n *gen.Type) string
	builder   Builder
	overwrite bool
}

func (w *writer) build(ctx context.Context, node *gen.Type, opts *BuildOptions) error {
	file, err := w.builder.Build(ctx, node, opts)
	if err != nil {
		return err
	}
	filename := w.filename(node)
	if !w.overwrite {
		// TODO
		if _, err = os.Stat(filename); err == nil || !os.IsNotExist(err) {
			return nil
		}
	}
	_ = os.MkdirAll(path.Dir(filename), 0700)
	return file.Save(filename)
}
