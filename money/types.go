package money

import "github.com/shopspring/decimal"

type Transaction struct {
	Date        string
	Description string
	Amount      decimal.Decimal
	Balance     decimal.Decimal
	Bank        string
	Account     string
	Tags        []string
	Category    Category
	Payee       string
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
	CategoryHouse         = "House"
	CategoryCar           = "Car"
	CategoryBills         = "Bills"
	CategoryIncome        = "Income"
	CategoryGroceries     = "Groceries"
	CategoryEntertainment = "Entertainment"
	CategoryVactation     = "Vacation"
	CategoryOther         = "Other"
)
