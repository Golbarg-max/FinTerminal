package cmd

import (

	"github.com/spf13/cobra"
	"fmt"
	"os"
)

var rootCmd = &cobra.Command {
	Use: "fin",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing request %s\n", err)
		os.Exit(1)
	}
}