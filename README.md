# FinTerminal 📊

A Bloomberg-inspired live market data terminal built in Go. Track your stock portfolio and view real-time market data from the command line.

## Features

✨ **Real-Time Stock Quotes** - Fetch current market data for any symbol
💼 **Portfolio Management** - Track holdings and calculate performance
📈 **Performance Metrics** - View gain/loss, ROI, and portfolio totals
🎨 **Bloomberg-Style Tables** - Professional formatted output with color coding
⚡ **Fast & Lightweight** - Native Go binary, minimal dependencies

## Quick Start

### 1. Get an API Key

Sign up for a free Alpha Vantage API key at [alphavantage.co](https://www.alphavantage.co/api/)

Free tier includes:
- 5 API calls per minute
- 500 calls per day
- Real-time stock data

### 2. Set Your API Key

**Option A: Using .env file (easiest)** ⭐

```bash
cp .env.example .env
# Edit .env and add your key
# ALPHA_VANTAGE_API_KEY="your_api_key_here"
```

**Option B: Environment variable**

```bash
export ALPHA_VANTAGE_API_KEY="your_api_key_here"
```

**Option C: Command flag**

```bash
fin quote AAPL --api-key "your_api_key_here"
```

See [ENV_SETUP.md](ENV_SETUP.md) for detailed configuration guide.

### 3. Build & Run

```bash
# Build the binary
go build -o fin

# Get a stock quote
./fin quote AAPL MSFT GOOGL

# View your portfolio
./fin portfolio
```

## Usage

### Get Stock Quotes

```bash
# Single symbol
fin quote AAPL

# Multiple symbols
fin quote AAPL MSFT GOOGL TSLA

# With custom API key
fin quote AAPL --api-key "your_key"
```

**Output:**
```
╔════════╦═════════╦════════════╦═══════════╦════════════╗
║ SYMBOL ║  PRICE  ║   CHANGE   ║ CHANGE %  ║  VOLUME    ║
╠════════╬═════════╬════════════╬═══════════╬════════════╣
║ AAPL   ║ $189.95 ║ +$2.50 (G) ║ +1.33% (G)║ 52,345,600 ║
╚════════╩═════════╩════════════╩═══════════╩════════════╝
```

### Portfolio Management

#### View Portfolio

```bash
fin portfolio
```

Displays your holdings with:
- Current price vs cost basis
- Total value and total invested
- Gain/Loss per position
- Return on Investment (ROI)
- Portfolio totals

**Output:**
```
💼 PORTFOLIO PERFORMANCE
════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════
║ SYMBOL ║ SHARES  ║ COST/SHARE ║ CURRENT ║ TOTAL COST ║ CURRENT VALUE ║ GAIN/LOSS ║ RETURN %  ║
╠════════╬═════════╬════════════╬═════════╬════════════╬═══════════════╬═══════════╬═══════════╣
║ AAPL   ║  10.00  ║   $150.00  ║ $189.95 ║  $1500.00  ║   $1899.50    ║ +$399.50  ║ +26.63%   ║
║ MSFT   ║  25.00  ║   $250.00  ║ $334.20 ║  $6250.00  ║   $8355.00    ║ +$2105.00 ║ +33.68%   ║
║ TSLA   ║   5.00  ║   $245.00  ║ $265.30 ║  $1225.00  ║   $1326.50    ║ +$101.50  ║ +8.29%    ║
╠════════╩═════════╩════════════╩═════════╩════════════╩═══════════════╩═══════════╩═══════════╣
║ TOTAL  ║         ║            ║         ║  $8975.00  ║  $11581.00    ║ +$2606.00 ║ +29.04%   ║
╚════════════════════════════════════════════════════════════════════════════════════════════════╝
```

#### Add Holdings

```bash
fin portfolio add AAPL 10 150.50
fin portfolio add MSFT 5 300.00
fin portfolio add GOOGL 2 2800.00
```

Adds shares with your cost basis (average cost per share)

#### Remove Holdings

```bash
fin portfolio remove AAPL
```

#### Specify Custom Portfolio File

```bash
fin portfolio --file /path/to/portfolio.json
fin portfolio add AAPL 10 150.50 --file /path/to/portfolio.json
```

By default, portfolio is saved to `~/.finterminal/portfolio.json`

## Portfolio File Format

Your portfolio is stored as JSON for easy backup and editing:

```json
{
  "holdings": {
    "AAPL": {
      "symbol": "AAPL",
      "quantity": 10,
      "cost_basis": 150.50,
      "current_price": 189.95,
      "last_updated": "2024-01-15T14:30:00Z"
    },
    "MSFT": {
      "symbol": "MSFT",
      "quantity": 25,
      "cost_basis": 250.00,
      "current_price": 334.20,
      "last_updated": "2024-01-15T14:30:00Z"
    }
  }
}
```

## Installation

### From Source

```bash
git clone https://github.com/yourusername/FinTerminal.git
cd FinTerminal
go build -o fin
sudo mv fin /usr/local/bin/  # Optional: add to PATH
```

### Prerequisites

- Go 1.21 or later
- Internet connection for API calls
- Free Alpha Vantage API key

## Architecture

```
FinTerminal/
├── main.go              # Entry point
├── cmd/
│   ├── root.go         # Root command
│   ├── quote.go        # Quote command
│   └── portfolio.go    # Portfolio commands
├── api/
│   └── alphavantage.go # API client
└── portfolio/
    ├── portfolio.go    # Portfolio data structure
    └── performance.go  # Performance calculations
```

## Dependencies

- [spf13/cobra](https://github.com/spf13/cobra) - CLI framework
- [olekukonko/tablewriter](https://github.com/olekukonko/tablewriter) - Table formatting

## Future Enhancements

- [ ] Support for multiple data providers (IEX, Finnhub, Polygon)
- [ ] Historical price charts (ASCII art)
- [ ] Portfolio alerts and notifications
- [ ] Dividend tracking
- [ ] Tax lot management
- [ ] Export to CSV/Excel
- [ ] Web dashboard
- [ ] Dark mode / custom themes
- [ ] Options data
- [ ] Cryptocurrency support

## Configuration

Set environment variables to customize behavior:

```bash
# API Configuration
export ALPHA_VANTAGE_API_KEY="your_key"

# Portfolio Configuration
export FINTERMINAL_PORTFOLIO_PATH="$HOME/.finterminal/portfolio.json"
```

## Troubleshooting

### "API key not provided"
```bash
export ALPHA_VANTAGE_API_KEY="your_api_key"
```

### "Symbol not found"
- Verify the symbol is correct (e.g., AAPL, not Apple)
- Check Alpha Vantage supported symbols
- Ensure your API key is valid

### "Rate limit exceeded"
- Alpha Vantage free tier: 5 calls/minute, 500 calls/day
- Upgrade to paid plan for higher limits
- Wait a minute before making more calls

### Portfolio file not found
- Portfolio is auto-created at `~/.finterminal/portfolio.json`
- Ensure directory exists or create manually
- Check file permissions

## Performance Notes

- First `portfolio` command takes 5-10 seconds (fetching data)
- Subsequent calls are instant
- API calls are cached in portfolio file
- Prices update when you run `portfolio` command

## Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

MIT License - see LICENSE file for details

## Author

Built with ❤️ by the FinTerminal community

## Support

- 📧 Email: support@finterminal.dev
- 🐛 Issues: github.com/yourusername/FinTerminal/issues
- 💬 Discussions: github.com/yourusername/FinTerminal/discussions

---

**Disclaimer:** This tool is for informational purposes. Not financial advice. Always do your own research before investing.
