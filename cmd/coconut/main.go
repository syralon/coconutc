package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/syralon/coconutc/internal/command"
)

const (
	usage = `A service generator base on ent(https://entgo.io/).
Homepage: https://github.com/syralon/coconutc.
`
)

func main() {
	cmd := &cobra.Command{
		Use:           "coconut",
		Long:          usage,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.AddCommand(
		command.Proto(),
		command.Service(),
		command.Generate(),
		command.Add(),
		command.New(),
	)
	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s", err.Error())
	}
}
