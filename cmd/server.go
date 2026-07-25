package cmd

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	handlers "FinTerminal/handlers"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the FinTerminal web server",
	Long:  "Start the FinTerminal web server with portfolio and quote endpoints",
	Run: func(cmd *cobra.Command, args []string) {
		// Load .env file
		if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
			log.Printf("Warning: Error loading .env file: %v\n", err)
		}

		// Get API key from flag or env
		if apiKey == "" {
			apiKey = os.Getenv("ALPHA_VANTAGE_API_KEY")
		}

		if apiKey == "" {
			fmt.Println("Error: API key not provided. Set ALPHA_VANTAGE_API_KEY environment variable or use --api-key flag")
			os.Exit(1)
		}

		port := cmd.Flag("port").Value.String()

		// Initialize handlers
		h := handlers.NewHandlers(apiKey)

		// Routes
		mux := http.NewServeMux()

		// Static files
		mux.Handle("/", http.FileServer(http.Dir("static")))

		// API endpoints
		mux.HandleFunc("/api/portfolio", h.HandlePortfolio)
		mux.HandleFunc("/api/holdings", h.HandleHoldings)
		mux.HandleFunc("/api/quote", h.HandleQuote)
		mux.HandleFunc("/api/prices/refresh", h.HandleRefreshPrices)

		// Start server
		addr := fmt.Sprintf(":%s", port)
		fmt.Printf("🚀 FinTerminal Web Server running on http://localhost%s\n", addr)
		log.Fatal(http.ListenAndServe(addr, mux))
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "Alpha Vantage API key")
	serverCmd.Flags().StringP("port", "p", "8080", "Port to run server on")
}