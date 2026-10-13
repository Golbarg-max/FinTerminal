package cmd

import (
	"FinTerminal/api"
	"FinTerminal/portfolio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var portfolioFilePath string

var portfolioCmd = &cobra.Command{
	Use:   "portfolio",
	Short: "Manage and display portfolio performance",
	Long:  "View your portfolio holdings and performance metrics",
	Run: func(cmd *cobra.Command, args []string) {
		// Use default path if not specified
		if portfolioFilePath == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Println("Error getting home directory:", err)
				os.Exit(1)
			}
			portfolioFilePath = filepath.Join(homeDir, ".finterminal", "portfolio.json")
		}

		// Load portfolio
		p := portfolio.NewPortfolio()
		if err := p.LoadFromFile(portfolioFilePath); err != nil {
			fmt.Println("Error loading portfolio:", err)
			os.Exit(1)
		}

		if len(p.Holdings) == 0 {
			fmt.Println("Portfolio is empty. Add holdings using 'portfolio add' command")
			return
		}

		client, _ := api.NewClient(apiKey)
		symbols := p.GetSymbols()

		fmt.Println("\n📊 Fetching latest prices...")
		for _, symbol := range symbols {
			quote, err := client.GetQuote(symbol)
			if err != nil {
				fmt.Printf("Error fetching %s: %v\n", symbol, err)
				continue
			}
			if err := p.UpdatePrice(symbol, quote.Price); err != nil {
				fmt.Printf("Warning: Failed to update price for %s: %v\n", symbol, err)
			}
		}

		// Display portfolio table
		displayPortfolioTable(p)

		// Save updated prices
		if err := p.SaveToFile(); err != nil {
			fmt.Println("Warning: Could not save portfolio:", err)
		}
	},
}

var portfolioAddCmd = &cobra.Command{
	Use:   "add SYMBOL QUANTITY COST_PER_SHARE",
	Short: "Add a holding to your portfolio",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		if portfolioFilePath == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Println("Error getting home directory:", err)
				os.Exit(1)
			}
			portfolioFilePath = filepath.Join(homeDir, ".finterminal", "portfolio.json")
		}

		p := portfolio.NewPortfolio()
		if err := p.LoadFromFile(portfolioFilePath); err != nil {
			fmt.Println("Error loading portfolio:", err)
			os.Exit(1)
		}

		symbol := args[0]
		quantity, err := parseFloat(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid quan %q: %v\n", args[1], err)
			os.Exit(1)
		}
		costBasis, err := parseFloat(args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid cost %q: %v\n", args[2], err)
			os.Exit(1)
		}

		if err := p.AddHolding(symbol, quantity, costBasis); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to add holding: %v\n", err)
			os.Exit(1)
		}
		if err := p.SaveToFile(); err != nil {
			fmt.Println("Error saving portfolio:", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Added %f shares of %s at $%.2f\n", quantity, symbol, costBasis)
	},
}

var portfolioRemoveCmd = &cobra.Command{
	Use:   "remove SYMBOL",
	Short: "Remove a holding from your portfolio",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if portfolioFilePath == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Println("Error getting home directory:", err)
				os.Exit(1)
			}
			portfolioFilePath = filepath.Join(homeDir, ".finterminal", "portfolio.json")
		}

		p := portfolio.NewPortfolio()
		if err := p.LoadFromFile(portfolioFilePath); err != nil {
			fmt.Println("Error loading portfolio:", err)
			os.Exit(1)
		}

		symbol := args[0]
		if err := p.RemoveHolding(symbol); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if err := p.SaveToFile(); err != nil {
			fmt.Println("Error saving portfolio:", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Removed %s from portfolio\n", symbol)
	},
}

func displayPortfolioTable(p *portfolio.Portfolio) {
	metrics := p.CalculatePerformance()
	totalCost, totalValue, totalGainLoss := p.CalculatePortfolioTotals()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"SYMBOL",
		"SHARES",
		"COST/SHARE",
		"CURRENT",
		"TOTAL COST",
		"CURRENT VALUE",
		"GAIN/LOSS",
		"RETURN %",
	})

	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_RIGHT)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT})

	// Add rows for each holding
	for _, m := range metrics {
		gainLossColor := ""
		returnColor := ""
		resetColor := "\033[0m"

		if m.GainLoss >= 0 {
			gainLossColor = "\033[32m" // Green
			returnColor = "\033[32m"
		} else {
			gainLossColor = "\033[31m" // Red
			returnColor = "\033[31m"
		}

		table.Append([]string{
			m.Symbol,
			fmt.Sprintf("%.2f", m.Quantity),
			fmt.Sprintf("$%.2f", m.CostPerShare),
			fmt.Sprintf("$%.2f", m.CurrentPrice),
			fmt.Sprintf("$%.2f", m.TotalCost),
			fmt.Sprintf("$%.2f", m.CurrentValue),
			fmt.Sprintf("%s$%.2f%s", gainLossColor, m.GainLoss, resetColor),
			fmt.Sprintf("%s%.2f%%%s", returnColor, m.ReturnOnInvestment, resetColor),
		})
	}

	// Add summary row
	summaryPercent := 0.0
	if totalCost > 0 {
		summaryPercent = (totalGainLoss / totalCost) * 100
	}

	summaryColor := ""
	if totalGainLoss >= 0 {
		summaryColor = "\033[32m"
	} else {
		summaryColor = "\033[31m"
	}
	resetColor := "\033[0m"

	table.Append([]string{
		"TOTAL",
		"",
		"",
		"",
		fmt.Sprintf("$%.2f", totalCost),
		fmt.Sprintf("$%.2f", totalValue),
		fmt.Sprintf("%s$%.2f%s", summaryColor, totalGainLoss, resetColor),
		fmt.Sprintf("%s%.2f%%%s", summaryColor, summaryPercent, resetColor),
	})

	fmt.Println("\n💼 PORTFOLIO PERFORMANCE")
	fmt.Println("=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=")
	table.Render()
	fmt.Println()
}

func parseFloat(s string) (float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	// ParseFloat accepts "NaN" and "Inf", but they aren't usable amounts:
	// NaN slips past `<= 0` checks because every comparison with it is false.
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, fmt.Errorf("must be a finite number")
	}
	return v, nil
}

func init() {
	rootCmd.AddCommand(portfolioCmd)
	portfolioCmd.AddCommand(portfolioAddCmd)
	portfolioCmd.AddCommand(portfolioRemoveCmd)

	portfolioCmd.Flags().StringVar(&portfolioFilePath, "file", "", "Path to portfolio file")
	portfolioCmd.Flags().StringVar(&apiKey, "api-key", "", "Alpha Vantage API key")
}
