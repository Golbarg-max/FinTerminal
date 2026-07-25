package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fin",
	Short: "FinTerminal - Bloomberg-inspired live market data terminal",
	Long: `FinTerminal is a command-line tool for tracking stock portfolios and viewing market data.

Get started:
  fin quote AAPL          - Get a quote for Apple
  fin portfolio           - View your portfolio performance
  fin portfolio add AAPL 10 150.50 - Add 10 shares of AAPL at $150.50`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
