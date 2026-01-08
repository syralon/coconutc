package command

import (
	"os/exec"
	"strings"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/spf13/cobra"
	"github.com/syralon/coconutc/internal/c/entc/entproto"
	"github.com/syralon/coconutc/internal/c/entc/entservice"
)

func Run() *cobra.Command {
	var opt options
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Generate api services from ent schemas.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if err := opt.parse(); err != nil {
				return err
			}
			cfg := &gen.Config{}
			graph, err := entc.LoadGraph(opt.target, cfg)
			if err != nil {
				return err
			}

			protoGenerator, err := entproto.NewGenerator(entproto.WithOutput(opt.output), entproto.WithVerbose(opt.verbose))
			if err != nil {
				return err
			}
			if err = protoGenerator.Generate(ctx, graph); err != nil {
				return err
			}
			
			serviceGenerator, err := entservice.NewGenerator(
				entservice.WithOutput(opt.output),
				entservice.WithOverwrite(opt.overwrite),
				entservice.WithEntPath(strings.TrimRight(opt.target, "/schema")),
			)
			if err != nil {
				return err
			}
			if err = serviceGenerator.Generate(cmd.Context(), graph); err != nil {
				return err
			}

			_ = exec.Command("buf", "format", "-w").Run()
			_ = exec.Command("ent", "generate", "--target", opt.target).Run()
			_ = exec.Command("go", "generate", "./...", opt.target).Run()
			return nil
		},
	}
	opt.register(cmd)
	return cmd
}
