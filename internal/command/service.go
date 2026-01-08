package command

import (
	"strings"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/spf13/cobra"
	"github.com/syralon/coconutc/internal/c/entc/entservice"
)

func Service() *cobra.Command {
	var opt options
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Generate api services from ent schemas.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opt.parse(); err != nil {
				return err
			}
			cfg := &gen.Config{}
			graph, err := entc.LoadGraph(opt.target, cfg)
			if err != nil {
				return err
			}
			serviceGenerator, err := entservice.NewGenerator(
				//entservice.WithModule(opt.module),
				entservice.WithOutput(opt.output),
				entservice.WithOverwrite(opt.overwrite),
				entservice.WithEntPath(strings.TrimRight(opt.target, "/schema")),
			)
			if err != nil {
				return err
			}
			return serviceGenerator.Generate(cmd.Context(), graph)
		},
	}
	opt.register(cmd)
	return cmd
}
