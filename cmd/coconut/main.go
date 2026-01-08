package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/syralon/coconutc/internal/command"
)

func main() {
	cmd := &cobra.Command{
		Use:           "entc-gen",
		Long:          "A service generator base on ent(https://entgo.io/).\nHomepage: https://github.com/syralon/coconutc",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.AddCommand(command.Ent()...)
	cmd.AddCommand(
		command.Proto(),
		command.Service(),
		command.Run(),
		command.New(),
	)
	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s", err.Error())
	}
}
