package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	handlers "FinTerminal/handlers"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the FinTerminal web server",
	Long:  "Start the FinTerminal web server with portfolio and quote endpoints",
	Run: func(cmd *cobra.Command, args []string) {

		// Precedence: --port flag > SERVER_PORT env > default.
		port, _ := cmd.Flags().GetString("port")
		if envPort := os.Getenv("SERVER_PORT"); envPort != "" && !cmd.Flags().Changed("port") {
			port = envPort
		}

		// Where the static frontend lives. Configurable so the binary can run
		// from a different working directory (e.g. inside a container).
		webDir := os.Getenv("FINTERMINAL_WEB_DIR")
		if webDir == "" {
			webDir = "web"
		}

		// Initialize handlers
		h := handlers.NewHandlers(apiKey)

		// Routes
		mux := http.NewServeMux()

		// Static files
		mux.Handle("/", http.FileServer(http.Dir(webDir)))

		// API endpoints
		mux.HandleFunc("/api/portfolio", h.HandlePortfolio)
		mux.HandleFunc("/api/holdings", h.HandleHoldings)
		mux.HandleFunc("/api/quote", h.HandleQuote)
		mux.HandleFunc("/api/prices/refresh", h.HandleRefreshPrices)
		mux.HandleFunc("/api/status", h.HandleStatus)

		srv := &http.Server{
			Addr:              fmt.Sprintf(":%s", port),
			Handler:           mux,
			ReadHeaderTimeout: 10 * time.Second,
		}

		fmt.Printf("🚀 FinTerminal Web Server running on http://localhost%s\n", srv.Addr)
		if h.IsDemoMode() {
			fmt.Printf("📊 Running in DEMO MODE (no API key needed)\n")
			fmt.Printf("💡 To use real data: Set ALPHA_VANTAGE_API_KEY in .env\n")
		} else {
			fmt.Printf("📊 API Key configured from environment\n")
		}
		fmt.Printf("💾 Portfolio file: %s\n", h.PortfolioPath())

		log.Fatal(srv.ListenAndServe())
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringP("port", "p", "8080", "Port to run server on")
}
