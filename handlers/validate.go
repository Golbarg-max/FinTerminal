package handlers

import (
	"regexp"

	"FinTerminal/api"
)

// A ticker is 1-10 letters/digits, optionally with '.' or '-' after the
// first character (BRK.B, BF-B, SHOP.TO). Compiled once at startup.
var symbolPattern = regexp.MustCompile(`^[A-Z0-9][A-Z0-9.-]{0,9}$`)

// validSymbol reports whether raw looks like a ticker once normalized.
// Checking here rejects junk with a 400 before it costs an API call.
func validSymbol(raw string) bool {
	return symbolPattern.MatchString(api.NormalizeSymbol(raw))
}
