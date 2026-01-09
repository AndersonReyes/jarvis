package money

import "regexp"

var incomeRegex, _ = regexp.Compile("payroll")

func income(tr *Transaction) bool {

	return incomeRegex.MatchString(tr.Description)
}
