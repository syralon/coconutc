package entservice

import (
	"context"
	"os"
	"testing"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/syralon/coconutc/pkg/annotation/entproto"
)

func TestRepositoryBuilder(t *testing.T) {
	cfg := &gen.Config{}
	graph, err := entc.LoadGraph("./testdata/ent/schema", cfg)
	if err != nil {
		t.Fatal(err)
	}

	var ctx = context.Background()
	builder := RepositoryBuilder("data", "github.com/syralon/example/internal/domain/entity", "github.com/syralon/example/internal/infra/tx")
	for _, node := range graph.Nodes {
		apiOpts, err := entproto.GetAPIOptions(node.Annotations)
		if err != nil {
			t.Error(err)
			return
		}
		opts := &BuildOptions{
			APIOptions:   &apiOpts,
			ProtoPackage: "github.com/syralon/example/proto/example",
			EntPackage:   "github.com/syralon/example/ent",
		}
		file, err := builder.Build(ctx, node, opts)
		if err != nil {
			t.Error(err)
			return
		}
		_ = file.Render(os.Stdout)
	}
}

func TestRepositoryInterfaceBuilder(t *testing.T) {
	cfg := &gen.Config{}
	graph, err := entc.LoadGraph("./testdata/ent/schema", cfg)
	if err != nil {
		t.Fatal(err)
	}

	var ctx = context.Background()
	builder := RepositoryInterfaceBuilder("repository", "github.com/syralon/example/internal/domain/entity")
	for _, node := range graph.Nodes {
		apiOpts, err := entproto.GetAPIOptions(node.Annotations)
		if err != nil {
			t.Error(err)
			return
		}
		opts := &BuildOptions{
			APIOptions:   &apiOpts,
			ProtoPackage: "github.com/syralon/example/proto/example",
			EntPackage:   "github.com/syralon/example/ent",
		}
		file, err := builder.Build(ctx, node, opts)
		if err != nil {
			t.Error(err)
			return
		}
		_ = file.Render(os.Stdout)
	}
}
