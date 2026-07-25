# FinTerminal - Project Structure Guide

## 📁 Directory Layout

```
FinTerminal/
├── main.go                      # Application entry point
├── go.mod                       # Go module definition
├── Makefile                     # Build automation
├── README.md                    # Project documentation
├── SETUP.md                     # Setup instructions
├── PROJECT_GUIDE.md             # This file
├── example-portfolio.json       # Example portfolio format
│
├── cmd/                         # CLI Commands
│   ├── root.go                  # Root command & help
│   ├── quote.go                 # Stock quote command
│   └── portfolio.go             # Portfolio commands (view/add/remove)
│
├── api/                         # External API Integration
│   └── alphavantage.go          # Alpha Vantage API client
│
└── portfolio/                   # Portfolio Logic
    ├── portfolio.go             # Portfolio data & file I/O
    └── performance.go           # Performance calculations

```

## 📄 File Descriptions

### Root Files

#### `main.go`
- **Purpose**: Application entry point
- **Responsibility**: Calls cmd.Execute() to start the CLI
- **Size**: ~7 lines
- **Key Functions**: main()

#### `go.mod`
- **Purpose**: Define project module and dependencies
- **Dependencies**:
  - `github.com/spf13/cobra` - CLI framework
  - `github.com/olekukonko/tablewriter` - Table formatting
- **Maintainers**: Go, when running `go get`

#### `Makefile`
- **Purpose**: Automate common tasks
- **Common targets**:
  - `make build` - Compile binary
  - `make run` - Build and run with args
  - `make install` - Install to /usr/local/bin
  - `make test` - Run tests
- **Time to Learn**: 5 minutes

### Documentation Files

#### `README.md`
- **Purpose**: Main project documentation
- **Covers**: Features, quick start, usage examples, troubleshooting
- **Audience**: End users
- **Size**: ~400 lines

#### `SETUP.md`
- **Purpose**: Step-by-step setup guide
- **Covers**: API key setup, building, testing, troubleshooting
- **Audience**: First-time users
- **Time to Complete**: ~8 minutes

#### `PROJECT_GUIDE.md`
- **Purpose**: Developer reference (this file)
- **Covers**: Architecture, code walkthrough, extending features
- **Audience**: Contributors, maintainers

### Configuration Files

#### `.gitignore`
- **Purpose**: Exclude files from git
- **Excludes**: Binaries, build artifacts, API keys, IDE files
- **Important**: Never commit `.finterminal/` directory or portfolio.json with real data

#### `example-portfolio.json`
- **Purpose**: Show users what portfolio.json looks like
- **Contains**: Sample holdings for AAPL, MSFT, GOOGL, TSLA
- **Use**: Copy and modify to create your own portfolio

### Code Modules

## 🏗️ Architecture Overview

### Layer 1: CLI (cmd/)

**Commands**:
- `fin quote AAPL` → Fetches and displays stock data
- `fin portfolio` → Shows performance metrics
- `fin portfolio add AAPL 10 150.50` → Adds holding
- `fin portfolio remove AAPL` → Removes holding

**Files**:
- `root.go` - Command tree root
- `quote.go` - Quote command implementation
- `portfolio.go` - Portfolio commands + table display

### Layer 2: API Integration (api/)

**Purpose**: Talk to external services (Alpha Vantage)

**Files**:
- `alphavantage.go`
  - `AlphaVantageClient` - API client struct
  - `GetQuote(symbol)` - Fetch live quote data
  - Returns: `QuoteData` struct with price, change, volume

**Key Struct**: `QuoteData`
```go
type QuoteData struct {
    Symbol        string    // "AAPL"
    Price         float64   // 189.95
    Change        float64   // 2.50
    ChangePercent float64   // 1.33
    Volume        int64     // 52345600
    Timestamp     time.Time
}
```

### Layer 3: Business Logic (portfolio/)

**Purpose**: Manage holdings and calculate performance

**Files**:
- `portfolio.go`
  - `Portfolio` - Main holdings container
  - `Holding` - Single stock position
  - Methods: Load/Save, Add/Remove, Get prices

- `performance.go`
  - `PerformanceMetrics` - Calculated metrics per holding
  - `CalculatePerformance()` - Compute all metrics
  - `CalculatePortfolioTotals()` - Sum across holdings

**Key Structs**:

`Portfolio`:
```go
type Portfolio struct {
    Holdings map[string]*Holding  // Symbol -> Holding
    filePath string                // Path to JSON file
}
```

`Holding`:
```go
type Holding struct {
    Symbol       string    // "AAPL"
    Quantity     float64   // 10.5 shares
    CostBasis    float64   // $150.50 per share
    CurrentPrice float64   // $189.95 (updated by API)
    LastUpdated  string    // "2024-01-15"
}
```

`PerformanceMetrics`:
```go
type PerformanceMetrics struct {
    Symbol              string
    Quantity            float64
    CostPerShare        float64
    CurrentPrice        float64
    TotalCost           float64    // Qty × CostBasis
    CurrentValue        float64    // Qty × CurrentPrice
    GainLoss            float64    // CurrentValue - TotalCost
    GainLossPercent     float64    // (GainLoss / TotalCost) × 100
    ReturnOnInvestment  float64
}
```

## 🔄 Data Flow

### Getting a Quote

```
User Input
    ↓
[quote.go] Parses args (AAPL, MSFT, ...)
    ↓
[alphavantage.go] Fetches from API
    ↓
[tablewriter] Formats output
    ↓
Display in Terminal
```

### Viewing Portfolio

```
User Input
    ↓
[portfolio.go] Load portfolio.json from disk
    ↓
[alphavantage.go] Fetch current prices for each holding
    ↓
[performance.go] Calculate metrics
    ↓
[portfolio.go] Build Bloomberg-style table
    ↓
[tablewriter] Format & display
    ↓
Save updated prices to portfolio.json
```

## 🔧 Key Functions Reference

### cmd/quote.go

```go
quoteCmd.Run()
├─ apiKey from flag or env
├─ NewAlphaVantageClient(apiKey)
├─ For each symbol:
│  ├─ client.GetQuote(symbol)
│  └─ Add row to table
└─ table.Render()
```

### cmd/portfolio.go

```go
portfolioCmd.Run()
├─ Load portfolio from file
├─ Get all symbols
├─ For each symbol:
│  ├─ Fetch current price
│  └─ Update in portfolio
├─ calculatePerformance()
├─ displayPortfolioTable()
└─ Save portfolio to file

portfolioAddCmd.Run()
├─ Load portfolio
├─ Parse: SYMBOL, QUANTITY, COST_PER_SHARE
├─ portfolio.AddHolding()
└─ Save portfolio

portfolioRemoveCmd.Run()
├─ Load portfolio
├─ portfolio.RemoveHolding(symbol)
└─ Save portfolio
```

### api/alphavantage.go

```go
GetQuote(symbol string) *QuoteData
├─ Build API URL with params
├─ HTTP GET request
├─ Parse JSON response
├─ Convert strings to numbers
└─ Return QuoteData struct
```

### portfolio/portfolio.go

```go
LoadFromFile(filepath)
├─ Check if file exists
├─ Read & parse JSON
└─ Populate Holdings map

SaveToFile()
├─ Create directory if needed
├─ Marshal Holdings to JSON
└─ Write to disk

AddHolding(symbol, qty, cost)
├─ Create Holding struct
└─ Add to Holdings map

UpdatePrice(symbol, price)
└─ Update CurrentPrice in holdings
```

### portfolio/performance.go

```go
CalculatePerformance() []PerformanceMetrics
├─ For each holding:
│  ├─ totalCost = qty × costBasis
│  ├─ currentValue = qty × currentPrice
│  ├─ gainLoss = currentValue - totalCost
│  ├─ gainLossPercent = (gainLoss / totalCost) × 100
│  └─ Create PerformanceMetrics
└─ Return slice of metrics

CalculatePortfolioTotals() (cost, value, gainLoss)
├─ Sum all totalCost
├─ Sum all currentValue
└─ Calculate portfolio gainLoss
```

## 🎯 Common Tasks

### Adding a New Data Provider

1. Create `api/yourprovider.go`
2. Implement `GetQuote(symbol)` method
3. Update `quote.go` to support --provider flag
4. Test with sample symbols

### Adding a New Command

1. Create `cmd/yourcommand.go`
2. Define `var yourCmd = &cobra.Command{ ... }`
3. In `init()`: `rootCmd.AddCommand(yourCmd)`
4. Add tests in `cmd/yourcommand_test.go`

### Extending Performance Metrics

1. Add fields to `PerformanceMetrics` struct
2. Calculate in `CalculatePerformance()`
3. Add column to `displayPortfolioTable()`

### Changing Table Format

1. Edit `displayPortfolioTable()` in `cmd/portfolio.go`
2. Modify `table.SetHeader()` for columns
3. Modify `table.Append()` for data rows
4. Use ANSI color codes for styling

## 📊 Portfolio File Location

By default: `~/.finterminal/portfolio.json`

Environment variable to override:
```bash
export FINTERMINAL_PORTFOLIO_PATH="/custom/path/portfolio.json"
```

## 🧪 Testing Checklist

- [ ] `fin quote AAPL` returns current price
- [ ] `fin quote AAPL MSFT GOOGL` handles multiple symbols
- [ ] `fin portfolio` shows empty message when no holdings
- [ ] `fin portfolio add AAPL 10 150.50` adds holding
- [ ] `fin portfolio` displays all holdings with performance
- [ ] `fin portfolio remove AAPL` removes holding
- [ ] Gain/Loss colors: GREEN for positive, RED for negative
- [ ] Portfolio saves/loads correctly to JSON
- [ ] Works with custom portfolio file path

## 🚀 Performance Notes

- **Quote API call**: 1-2 seconds per symbol
- **Portfolio load**: <100ms (local file read)
- **Portfolio display**: 5-10 seconds (API calls for all holdings)
- **File I/O**: <50ms for typical portfolio

## 🔐 Security Notes

⚠️ **Never commit your API key!**
- Use environment variables
- Use .gitignore to exclude portfolio.json
- Don't share portfolio files with real data

## 📚 Dependencies

### Cobra (CLI Framework)
- Purpose: Command parsing, help text, subcommands
- Why: Industry standard Go CLI framework
- Learning: 30 minutes
- Docs: https://github.com/spf13/cobra

### TableWriter (Table Formatting)
- Purpose: Format data as ASCII tables
- Why: Professional-looking terminal output
- Learning: 15 minutes
- Docs: https://github.com/olekukonko/tablewriter

## 🔄 Update Workflow

```bash
# 1. Make changes to files
# 2. Run tests
make test

# 3. Build locally
make build

# 4. Test commands
./fin quote AAPL

# 5. Commit
git add .
git commit -m "Add feature"

# 6. Push
git push origin feature-branch
```

---

**Next Steps**: Start exploring the code! Pick a file and read through it. Try making small modifications. Build and test.
