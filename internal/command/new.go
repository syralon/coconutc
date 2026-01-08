package command

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/spf13/cobra"
	"github.com/syralon/coconutc/internal/c/entc/entschema"
)

const (
	newUsage = `Create new ent schemas. Usage: coconut new <SchemaName>([field_1],[field_2]:[field_type]), eg:
	coconut new User
	coconut new User(id)
	coconut new User(id:int64, firstname, lastname, email)
`
)

func New() *cobra.Command {
	var opt options
	cmd := &cobra.Command{
		Use:   "new",
		Short: newUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opt.parse(); err != nil {
				return err
			}
			cfg := &gen.Config{}
			graph, err := entc.LoadGraph(opt.target, cfg)
			if err != nil {
				return err
			}
			return entschema.Generate(opt.target, args, graph)
		},
	}
	opt.register(cmd)
	return cmd
}
