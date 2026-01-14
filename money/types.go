package money

import (
	"github.com/shopspring/decimal"
)

var AccountIdToFullName = map[string]string{
	"7150": "Discover Checking 7150",
	"0075": "Discover Savings 0075",
	"5029": "Discover Kids Savings 5029",
	"8444": "Capital one credit card 8444",
}

type Transaction struct {
	Date               string
	Description        string
	Amount             decimal.Decimal
	Balance            decimal.Decimal
	Bank               string
	SourceAccount      string
	DestinationAccount string
	Tags               []string
	Category           Category
	Payee              string
}

var TransactionColumns = [7]string{"Date", "SourceAccount", "Category", "Amount",
	"Description", "Payee", "DestinationAccount"}

func TransactionColumnsToIndex() map[string]int {
	m := make(map[string]int)
	for idx, col := range TransactionColumns {
		m[col] = idx
	}

	return m
}

func NewTransactionFromValues(colValues []string) Transaction {
	mapping := TransactionColumnsToIndex()
	return Transaction{
		Date:               colValues[mapping["Date"]],
		Description:        colValues[mapping["Description"]],
		Amount:             decimal.RequireFromString(colValues[mapping["Amount"]]),
		SourceAccount:      colValues[mapping["SourceAccount"]],
		DestinationAccount: colValues[mapping["DestinationAccount"]],
		Category:           colValues[mapping["Category"]],
	}
}

//type BillId = int
//
//const (
//	BillIdMortgage BillId = iota
//	BillIdWaterAndSewerm
//	BillIdInternet
//	BillIdPhoneBill
//	BillIdGas
//	BillIdCarInsurance
//	BillIdHomeInsuraznce
//	BillIdCorollaNote
//	BillIdDogInsurance
//	BillIdXanderNinjaWarrior
//	BillIdXanderPiano
//	BillIdNetflix
//	BillIdYoutubePremium
//	BillIdElectric
//	BillIdEzPass
//	BillIdCkoFitness
//	BillIdCassBudget
//	BillIdCassClubPilates
//)

type Category = string

const (
	CategoryHouse             = "House"
	CategoryCar               = "Car"
	CategoryBills             = "Bills"
	CategoryIncome            = "Income"
	CategoryGroceries         = "Groceries"
	CategoryEntertainment     = "Entertainment"
	CategoryVactation         = "Vacation"
	CategoryOther             = "Other"
	CategoryCreditCardPayment = "Payment/Credit"
	CategoryTransfer          = "Transfer"
	CategoryUnknown           = "UNKNOWN"
)
