package cmd

import (
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

		apiKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
		if apiKey == "" {
			log.Fatal("demo")
		}

		port := os.Getenv("SERVER_PORT")
		if port == "" {
			port = "8080"
		}

		// Initialize handlers
		h := handlers.NewHandlers(apiKey)

		// Routes
		mux := http.NewServeMux()

		// Static files
		mux.Handle("/", http.FileServer(http.Dir("web")))

		// API endpoints
		mux.HandleFunc("/api/portfolio", h.HandlePortfolio)
		mux.HandleFunc("/api/holdings", h.HandleHoldings)
		mux.HandleFunc("/api/quote", h.HandleQuote)
		mux.HandleFunc("/api/prices/refresh", h.HandleRefreshPrices)
		mux.HandleFunc("/api/status", h.HandleStatus)

		// Start server
		addr := fmt.Sprintf(":%s", port)
		fmt.Printf("🚀 FinTerminal Web Server running on http://localhost%s\n", addr)

		if apiKey == "demo" {
			fmt.Printf("📊 Running in DEMO MODE (no API key needed)\n")
			fmt.Printf("💡 To use real data: Set ALPHA_VANTAGE_API_KEY in .env\n")
		} else {
			fmt.Printf("📊 API Key configured from environment\n")
		}

		log.Fatal(http.ListenAndServe(addr, mux))
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringP("port", "p", "8080", "Port to run server on")
}