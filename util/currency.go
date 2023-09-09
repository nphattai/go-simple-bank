package util

import "golang.org/x/exp/slices"

var SupportedCurrency = []string{"USD", "EUR", "CAD"}

func IsSupportedCurrency(currency string) bool {
	return slices.Contains(SupportedCurrency, currency)
}
