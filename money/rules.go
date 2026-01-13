package money

import (
	"regexp"
	"strings"
)

// TODO: explore having array of regexes map to a category instead one big regex
//var bills = []string{
//	".*verizon.*",
//}

var billsRegex = regexp.MustCompile("(verizon|cloudflare|progressive|bank of america|head over heels|clubpilate|pets best ins|youtubepremium|netflix|pseg|bsi financial|cko|octopus music school)")
var incomeRegex = regexp.MustCompile("(payroll ach from spotify|ach deposit|interest paid|cashback)")
var groceriesRegex = regexp.MustCompile("(walmart|bjs|walgreens|target|wal-mart)")
var creditCardPaymentRegex = regexp.MustCompile("(capital one online pmt)")
var transferCategoryRegex = regexp.MustCompile("(transfer to|transfer from)")

// items that dont get categorized by upstream banks
var uncategorizedRegex = regexp.MustCompile("(zelle)")

func SetCategory(bankCategory string, transaction *Transaction) {

	var cat = bankCategory
	switch bankCategory {
	case "Gas/Automotive":
		cat = CategoryCar
	case "Other Services":
		cat = CategoryOther
	case "Lodging":
		cat = CategoryVactation
	default:
		desc := strings.ToLower(transaction.Description)
		if billsRegex.MatchString(desc) {
			cat = CategoryBills
		} else if incomeRegex.MatchString(desc) {
			cat = CategoryIncome
		} else if groceriesRegex.MatchString(desc) {
			cat = CategoryGroceries
		} else if uncategorizedRegex.MatchString(desc) {
			cat = CategoryOther
		} else if creditCardPaymentRegex.MatchString(desc) {
			cat = CategoryCreditCardPayment
		} else if transferCategoryRegex.MatchString(desc) {
			cat = CategoryTransfer
		}
	}

	transaction.Category = cat
}

var skipPayeeRegex = regexp.MustCompile("(transfer to|transfer from|deposit from|interest paid)")

func SetPayee(transaction *Transaction) {
	desc := strings.ToLower(transaction.Description)

	// only set the payee for bills so we can auto track with schedules
	if skipPayeeRegex.MatchString(desc) || transaction.Category != CategoryBills {
		return
	}

	var payee = ""

	var regexes = []regexp.Regexp{}
	regexes = append(regexes, *billsRegex)
	regexes = append(regexes, *incomeRegex)
	regexes = append(regexes, *groceriesRegex)
	regexes = append(regexes, *uncategorizedRegex)

	for _, r := range regexes {

		var match = r.FindStringIndex(desc)

		if match != nil {
			payee = desc[match[0]:match[1]]
		}
	}

	transaction.Payee = payee
}
