package entservice

import (
	"context"
	"testing"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func TestGenerator(t *testing.T) {
	cfg := &gen.Config{}
	graph, err := entc.LoadGraph("./testdata/ent/schema", cfg)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	g, err := NewGenerator(WithOutput("../../../../example"), WithOverwrite(true))
	if err != nil {
		t.Fatal(err)
	}
	if err = g.Generate(ctx, graph); err != nil {
		t.Fatal(err)
	}
}
