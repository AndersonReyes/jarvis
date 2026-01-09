package money

import "github.com/shopspring/decimal"

type Transaction struct {
	Date        string
	Description string
	Amount      decimal.Decimal
	Balance     decimal.Decimal
	Category    string
	Account     string
	Tags        []string
}
