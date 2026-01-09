package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var opt options
	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new project",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("usage: coconut new Example")
			}
			module := args[0]
			name := module
			if n := strings.LastIndex(module, "/"); n > 0 {
				name = module[n+1:]
			}
			if _, err := os.Stat(name); err == nil {
				return fmt.Errorf("%s already exists", name)
			}
			if err := os.Mkdir(name, os.ModePerm); err != nil {
				return err
			}
			if err := runExecCommands(
				name,
				opt.verbose,
				exec.CommandContext(cmd.Context(), "go", "mod", "init", module),
				exec.CommandContext(cmd.Context(), "go", "get", "entgo.io/ent"),
				exec.CommandContext(cmd.Context(), "coconut", "add", "Example(name,status:int)"),
				exec.CommandContext(cmd.Context(), "coconut", "generate"),
				exec.CommandContext(cmd.Context(), "ent", "generate", opt.target),
			); err != nil {
				return err
			}
			fmt.Printf("Project '%s' is created.\n", module)
			fmt.Printf("Quick start: 'cd %s && go run'\n", name)
			return nil
		},
	}
	opt.register(cmd)
	return cmd
}

func runExecCommands(dir string, verbose bool, commands ...*exec.Cmd) error {
	for _, cmd := range commands {
		cmd.Dir = dir
		cmd.Stderr = os.Stderr
		if verbose {
			fmt.Println(cmd.String())
		}
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
