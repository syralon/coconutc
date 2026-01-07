package entservice

import (
	"context"
	"os"
	"testing"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/syralon/coconutc/pkg/annotation/entproto"
)

func TestServiceBuilder(t *testing.T) {
	cfg := &gen.Config{}
	graph, err := entc.LoadGraph("./testdata/ent/schema", cfg)
	if err != nil {
		t.Fatal(err)
	}

	var ctx = context.Background()
	builder := ServiceBuilder("service", "github.com/syralon/example/internal/repository")
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
