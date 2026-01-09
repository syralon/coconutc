package command

import (
	"os"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/spf13/cobra"
	"github.com/syralon/coconutc/internal/c/entc/entschema"
)

const (
	addUsage = `"Add or modify ent schemas. 
Add a '-' after name to overwrite existed fields; 
If you not specified field type, the field name end with '_id' will be treated as int, and the field name end with '_at' will be treated as Time.
Usage: coconut add <SchemaName>([field_1],[field_2]:[field_type]), eg:
	coconut add User
	coconut add 'User(id)'
	coconut add 'User(id:int64, firstname, lastname, email)'`
)

func Add() *cobra.Command {
	var opt options
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add or modify ent schemas. run `coconut help add` for more info.",
		Long:  addUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opt.parse(); err != nil {
				return err
			}
			cfg := &gen.Config{}

			var graph *gen.Graph
			if _, err := os.Stat(opt.target); !os.IsNotExist(err) {
				graph, err = entc.LoadGraph(opt.target, cfg)
				if err != nil {
					return err
				}
			}
			return entschema.Generate(opt.target, args, graph)
		},
	}
	opt.register(cmd)
	return cmd
}
