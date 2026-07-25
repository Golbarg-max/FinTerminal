# FinTerminal Setup Guide

## Step 1: Get an API Key (2 minutes)

1. Visit [Alpha Vantage](https://www.alphavantage.co/api/)
2. Enter your email
3. You'll receive your API key immediately
4. No credit card required for free tier

## Step 2: Clone & Build (3 minutes)

```bash
# Clone the repository
git clone https://github.com/yourusername/FinTerminal.git
cd FinTerminal

# Build the binary
make build

# Or manually:
go build -o fin
```

## Step 3: Set Your API Key (1 minute)

### Option A: Environment Variable (Recommended)

```bash
cp .env.example .env
# Edit .env and add your API key
# ALPHA_VANTAGE_API_KEY="your_api_key_here"

See [ENV_SETUP.md](ENV_SETUP.md) for detailed configuration guide.
```

### Option B: Command Line Flag
```bash
export ALPHA_VANTAGE_API_KEY="your_api_key_here"
source ~/.bashrc
```
### Option C: Command Line Flag

```bash
fin quote AAPL --api-key "your_api_key_here"
```
## Step 4: Test Installation (1 minute)

```bash
# Get a stock quote
fin quote AAPL

# You should see:
# ╔════════╦═════════╦════════════╦═══════════╦════════════╗
# ║ SYMBOL ║  PRICE  ║   CHANGE   ║ CHANGE %  ║  VOLUME    ║
# ╠════════╬═════════╬════════════╬═══════════╬════════════╣
# ║ AAPL   ║ $189.95 ║ +$2.50     ║ +1.33%    ║ 52,345,600 ║
# ╚════════╩═════════╩════════════╩═══════════╩════════════╝
```

## Step 5: Create Your Portfolio (2 minutes)

```bash
# Add some holdings
fin portfolio add AAPL 10 150.50
fin portfolio add MSFT 5 300.00
fin portfolio add GOOGL 2 2800.00

# View your portfolio performance
fin portfolio
```

## Troubleshooting

### "Command not found: fin"

**Option 1: Add to PATH**
```bash
# Move binary to /usr/local/bin
sudo make install

# Or manually:
sudo cp fin /usr/local/bin/
```

**Option 2: Run from project directory**
```bash
./fin quote AAPL
```

### "API key not provided"

Make sure your environment variable is set:
```bash
# Check if it's set
echo $ALPHA_VANTAGE_API_KEY

# If empty, add to your shell profile:
export ALPHA_VANTAGE_API_KEY="your_key"
```

### "Symbol not found"

- Double-check the ticker symbol (e.g., AAPL, not Apple)
- Alpha Vantage may not have all symbols
- Check [this list](https://www.alphavantage.co/) for supported symbols

### Rate Limit Issues

Free Alpha Vantage includes:
- 5 API calls per minute
- 500 calls per day

Wait a minute between commands or [upgrade to paid plan](https://www.alphavantage.co/premium/)

## Next Steps

1. **Add more holdings** to your portfolio
2. **Set up daily monitoring** with a cron job:
   ```bash
   # Check portfolio every morning at 9 AM
   0 9 * * 1-5 fin portfolio >> ~/portfolio.log
   ```
3. **Export data** for further analysis
4. **Contribute** improvements back to the project

## Getting Help

- 📖 Read the [README.md](README.md)
- 🐛 Check [Issues](https://github.com/yourusername/FinTerminal/issues)
- 💬 Start a [Discussion](https://github.com/yourusername/FinTerminal/discussions)

## What's Next?

Your portfolio is saved at: `~/.finterminal/portfolio.json`

You can:
- Edit it manually (it's just JSON)
- Back it up to cloud storage
- Share holdings list with others (but keep API key secret!)

---

**Happy tracking! 📊**
