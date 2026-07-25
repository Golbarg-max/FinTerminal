# FinTerminal - Quick Reference

## 🚀 Quick Start

```bash
# Set API key
export ALPHA_VANTAGE_API_KEY="your_key"

# Build
make build

# Get quotes
fin quote AAPL MSFT

# View portfolio
fin portfolio

# Manage portfolio
fin portfolio add AAPL 10 150.50
fin portfolio remove AAPL
```

## 📖 Commands

### Quote Command
```bash
fin quote AAPL                              # Single stock
fin quote AAPL MSFT GOOGL TSLA              # Multiple stocks
fin quote AAPL --api-key "key_here"         # Custom API key
```

**Output Columns**: SYMBOL | PRICE | CHANGE | CHANGE % | VOLUME

### Portfolio View
```bash
fin portfolio                               # Show all holdings
fin portfolio --file /path/to/portfolio.json # Custom file
fin portfolio --api-key "key_here"          # Custom API key
```

**Output Columns**: 
- SYMBOL | SHARES | COST/SHARE | CURRENT
- TOTAL COST | CURRENT VALUE | GAIN/LOSS | RETURN %

### Portfolio Add
```bash
fin portfolio add AAPL 10 150.50            # Add holding
fin portfolio add MSFT 25 300.00 --file ~/my-portfolio.json
```

**Format**: `fin portfolio add SYMBOL QUANTITY COST_PER_SHARE`

### Portfolio Remove
```bash
fin portfolio remove AAPL                   # Remove holding
fin portfolio remove MSFT --file ~/my-portfolio.json
```

## 🔑 API Setup

1. Get key: https://www.alphavantage.co/api/
2. Set environment: `export ALPHA_VANTAGE_API_KEY="your_key"`
3. Or use flag: `--api-key "your_key"`

**Limits (Free Tier)**:
- 5 calls/minute
- 500 calls/day
- Real-time data

## 📊 Portfolio JSON Format

```json
{
  "holdings": {
    "AAPL": {
      "symbol": "AAPL",
      "quantity": 10,
      "cost_basis": 150.50,
      "current_price": 189.95,
      "last_updated": "2024-01-15"
    }
  }
}
```

**Location**: `~/.finterminal/portfolio.json`

## 🎯 Performance Calculations

```
Total Cost    = Quantity × Cost Per Share
Current Value = Quantity × Current Price
Gain/Loss     = Current Value - Total Cost
Return %      = (Gain/Loss / Total Cost) × 100
```

**Colors**:
- 🟢 **Green**: Positive gain/return
- 🔴 **Red**: Negative loss/return

## 🛠️ Make Commands

```bash
make build      # Compile binary
make run        # Build & run
make test       # Run tests
make clean      # Remove binaries
make install    # Install to /usr/local/bin
make deps       # Download dependencies
make fmt        # Format code
make quote      # Test quote command
make portfolio  # Test portfolio command
make example    # Add example holdings
```

## 🌍 Environment Variables

```bash
ALPHA_VANTAGE_API_KEY     # Required: API key
FINTERMINAL_PORTFOLIO_PATH # Optional: Custom portfolio file path
```

## 📂 File Structure

```
FinTerminal/
├── cmd/               # CLI commands
│   ├── root.go       # Root command
│   ├── quote.go      # Stock quotes
│   └── portfolio.go  # Portfolio management
├── api/              # API clients
│   └── alphavantage.go
├── portfolio/        # Business logic
│   ├── portfolio.go   # Data structure
│   └── performance.go # Calculations
├── main.go           # Entry point
├── go.mod            # Dependencies
├── Makefile          # Build automation
└── README.md         # Documentation
```

## 🆘 Troubleshooting

| Problem | Solution |
|---------|----------|
| "Command not found: fin" | Run `make install` or `./fin` from project dir |
| "API key not provided" | Set `ALPHA_VANTAGE_API_KEY` environment variable |
| "Symbol not found" | Check ticker symbol (e.g., AAPL not Apple) |
| "Rate limit exceeded" | Free tier: 5 calls/min, 500 calls/day |
| "Portfolio not found" | Auto-created at `~/.finterminal/portfolio.json` |

## 🔄 Typical Workflow

```bash
# 1. Initial setup
export ALPHA_VANTAGE_API_KEY="your_key_here"
make build

# 2. Check a stock
fin quote AAPL

# 3. Start portfolio
fin portfolio add AAPL 10 150.50

# 4. Add more holdings
fin portfolio add MSFT 25 300.00
fin portfolio add GOOGL 5 2800.00

# 5. Monitor performance
fin portfolio

# 6. Update holdings
fin portfolio remove AAPL
fin portfolio add AAPL 15 155.00  # New position
```

## 📈 Real-World Examples

```bash
# Day trader monitoring
fin quote AAPL MSFT GOOGL TSLA NVDA AMZN

# Compare your holdings
fin portfolio

# Check single stock fundamentals
fin quote BRK.B

# Sector comparison
fin quote XLK XLV XLF XLE  # Tech, Health, Finance, Energy

# Tech portfolio snapshot
fin portfolio add META 5 200.00
fin portfolio add NVDA 3 850.00
fin portfolio add AAPL 10 150.00
fin portfolio  # See performance
```

## ⚙️ Configuration Tips

**Daily monitoring** (add to crontab):
```bash
0 9 * * 1-5 /usr/local/bin/fin portfolio >> ~/portfolio.log
```

**Multiple portfolios**:
```bash
fin portfolio --file ~/portfolios/tech.json
fin portfolio --file ~/portfolios/dividend.json
```

**Batch quotes**:
```bash
for stock in AAPL MSFT GOOGL TSLA NVDA; do
  fin quote $stock
done
```

## 🎨 Display Tips

- Prices formatted as currency: $189.95
- Percentages with % sign: +1.33%
- Large numbers with commas: 52,345,600
- Green for gains, Red for losses
- Tables with borders for clarity

## 💡 Performance Notes

- First portfolio load: 5-10 seconds (API calls)
- Subsequent loads: <1 second (cached prices)
- Quote for single symbol: 1-2 seconds
- File operations: <50ms

---

**Pro Tips**:
- ✅ Check portfolio daily for market changes
- ✅ Update cost basis when adding to positions
- ✅ Remove positions when you sell
- ✅ Back up portfolio.json to cloud storage
- ✅ Use for tracking only (not financial advice!)
