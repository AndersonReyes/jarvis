package money

import "github.com/shopspring/decimal"

var Budget = map[Category]decimal.Decimal{
	CategoryBills:     decimal.NewFromFloat(4638.08),
	CategoryGroceries: decimal.NewFromFloat(1000.0),
}
