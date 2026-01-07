package main

import (
	"context"
	"flag"
	"os/exec"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/syralon/coconutc/internal/c/entc/entproto"
)

var target string
var verbose bool

func init() {
	flag.StringVar(&target, "target", "./ent/schema", "")
	flag.BoolVar(&verbose, "v", false, "")
	flag.Parse()
}

func main() {
	ctx := context.Background()
	cfg := &gen.Config{}
	graph, err := entc.LoadGraph(target, cfg)
	if err != nil {
		panic(err)
	}
	generator, err := entproto.NewGenerator(entproto.WithOutput("."), entproto.WithVerbose(verbose))
	if err != nil {
		panic(err)
	}
	if err = generator.Generate(ctx, graph); err != nil {
		panic(err)
	}
	_ = exec.Command("buf", "format", "-w").Run()
}
