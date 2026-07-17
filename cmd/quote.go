package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
)

var quoteCmd = &cobra.Command{
	Use: "quote",
	Short: "Print a quote for one or more symbols",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, v := range args {
			fmt.Println(v )
		}
	},
}

func init() {
	rootCmd.AddCommand(quoteCmd)
}

