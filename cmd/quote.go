package cmd

import (
	"FinTerminal/api"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var apiKey string

var quoteCmd = &cobra.Command{
	Use:   "quote [SYMBOLS...]",
	Short: "Fetch real-time stock quotes",
	Long:  "Display current market data for one or more stock symbols",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		
		client, _ := api.NewClient(apiKey)


		// Create table
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"SYMBOL", "PRICE", "CHANGE", "CHANGE %", "VOLUME"})
		table.SetBorder(true)
		table.SetRowLine(true)
		table.SetAlignment(tablewriter.ALIGN_RIGHT)

		// Fetch quotes for each symbol
		for _, symbol := range args {
			quote, err := client.GetQuote(symbol)
			if err != nil {
				fmt.Printf("Error fetching %s: %v\n", symbol, err)
				continue
			}

			changeColor := ""
			changePercentColor := ""
			if quote.Change >= 0 {
				changeColor = "\033[32m" // Green
				changePercentColor = "\033[32m"
			} else {
				changeColor = "\033[31m" // Red
				changePercentColor = "\033[31m"
			}
			resetColor := "\033[0m"

			table.Append([]string{
				quote.Symbol,
				fmt.Sprintf("$%.2f", quote.Price),
				fmt.Sprintf("%s$%.2f%s", changeColor, quote.Change, resetColor),
				fmt.Sprintf("%s%.2f%%%s", changePercentColor, quote.ChangePercent, resetColor),
				fmt.Sprintf("%d", quote.Volume),
			})
		}

		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(quoteCmd)
	quoteCmd.Flags().StringVar(&apiKey, "api-key", "", "Alpha Vantage API key")
}
