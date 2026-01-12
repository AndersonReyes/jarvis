package money

import (
	"regexp"
	"strings"
)

// TODO: explore having array of regexes map to a category instead one big regex
//var bills = []string{
//	".*verizon.*",
//}

var billsRegex = regexp.MustCompile(".*(verizon|cloudflare|progressive|bank of america|head over heels|clubpilate|pets|youtubepremium|netflix|pseg|bsi financial).*")
var incomeRegex = regexp.MustCompile(".*(payroll ach from Spotify|ach deposit|interest paid|cashback).*")
var groceriesRegex = regexp.MustCompile(".*(walmart|bjs|walgreens|target).*")

// items that dont get categorized by upstream banks
var uncategorizedRegex = regexp.MustCompile(".*(zelle).*")

func GetCategory(bankCategory string, transaction Transaction) Category {

	switch bankCategory {
	case "Gas/Automotive":
		return CategoryCar
	case "Other Services":
		return CategoryOther
	case "Lodging":
		return CategoryVactation
	default:
		desc := strings.ToLower(transaction.Description)
		if billsRegex.MatchString(desc) {
			return CategoryBills
		} else if incomeRegex.MatchString(desc) {
			return CategoryIncome
		} else if groceriesRegex.MatchString(desc) {
			return CategoryGroceries
		} else if uncategorizedRegex.MatchString(desc) {
			return CategoryOther
		}
		return bankCategory

	}
}

func GetPayee(account string, transaction Transaction) string {
	desc := strings.ToLower(transaction.Description)

	var regexes = []regexp.Regexp{}
	regexes = append(regexes, *billsRegex)
	regexes = append(regexes, *incomeRegex)
	regexes = append(regexes, *groceriesRegex)
	regexes = append(regexes, *uncategorizedRegex)

	for _, r := range regexes {

		var match = r.FindStringIndex(desc)

		if match != nil {
			// get 50 more chars past the end as payee
			end := match[1] + 50
			return desc[match[0]:end]
		}
	}

	return "(no Name)"
}
