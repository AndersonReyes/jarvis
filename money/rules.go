package money

import (
	"log"
	"regexp"
	"strings"
)

var billsRegex = regexp.MustCompile("(verizon|cloudflare|progressive|bank of america|head over heels|clubpilate|pets best ins|youtubepremium|netflix|pseg|bsi financial|cko|octopus music school)")
var incomeRegex = regexp.MustCompile("(payroll ach from spotify|ach deposit|interest paid|cashback)")
var groceriesRegex = regexp.MustCompile("(walmart|bjs|walgreens|target|wal-mart)")
var creditCardPaymentRegex = regexp.MustCompile("(capital one online pmt)")
var transferCategoryRegex = regexp.MustCompile("(transfer to|transfer from)")

// items that dont get categorized by upstream banks
var otherCategoryRegex = regexp.MustCompile("(zelle)")

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
		} else if otherCategoryRegex.MatchString(desc) {
			cat = CategoryOther
		} else if creditCardPaymentRegex.MatchString(desc) {
			cat = CategoryCreditCardPayment
		} else if transferCategoryRegex.MatchString(desc) {
			cat = CategoryTransfer
		} else if cat == "" {
			cat = CategoryUnknown
		}
	}

	transaction.Category = cat
}

var skipPayeeRegex = regexp.MustCompile("(transfer from|deposit from|interest paid)")

func SetPayeeAndDestinationAccount(transaction *Transaction) {
	desc := strings.ToLower(transaction.Description)

	// only set the payee for bills so we can auto track with schedules
	if skipPayeeRegex.MatchString(desc) {
		return
	}

	if transaction.Category == CategoryTransfer {
		externalAcc := AccountIdToFullName[desc[len(desc)-4:]]
		transaction.DestinationAccount = externalAcc
		log.Printf("Processed transfer %+v", transaction)
	}

	var payee = ""

	var regexes = []regexp.Regexp{}
	regexes = append(regexes, *billsRegex)
	regexes = append(regexes, *incomeRegex)
	regexes = append(regexes, *groceriesRegex)
	regexes = append(regexes, *otherCategoryRegex)

	for _, r := range regexes {

		var match = r.FindStringIndex(desc)

		if match != nil {
			payee = desc[match[0]:match[1]]
		}
	}

	transaction.Payee = payee
}
