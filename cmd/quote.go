package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
)

var quoteCmd = &cobra.Command{
	Use: "quote",
	Short: "",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fin" + args[0] + args[1])
	},
}

func init() {
	rootCmd.AddCommand(quoteCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing quote %s\n", err)
		os.Exit(1)
	}
}