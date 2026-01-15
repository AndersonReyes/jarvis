package money

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

var FireflyCmd = &cobra.Command{
	Use:   "firefly [csv]",
	Short: "deploy csv to firefly",
	Long:  "create the transactions in one shot",
	Args:  cobra.MinimumNArgs(1),
	Run:   fireflyRun,
}

var (
	importTransactionsFlag bool
	allowDuplicates        bool
)

var client = NewFireflyClient()

func init() {
	FireflyCmd.Flags().BoolVarP(&importTransactionsFlag, "transactions", "t",
		false, "import transactions from csv")
	FireflyCmd.Flags().BoolVarP(&allowDuplicates, "duplicates", "d",
		false, "allow duplicates on upload")
}

func fireflyRun(cmd *cobra.Command, args []string) {
	log.SetPrefix("money.firefly3: ")

	if importTransactionsFlag {
		importTransactions(args)
	} else {
		log.Fatalf("Need a flag to run: jarvis firefly --help")
	}

}

func getType(t Transaction) string {
	if t.Category == CategoryTransfer {
		return "transfer"
	} else if t.Amount.LessThan(decimal.NewFromFloat(0.0)) {
		return "withdrawal"
	}

	return "deposit"
}

func importTransactions(args []string) {
	filename := args[0]
	f, err := os.Open(filename)
	defer func() {
		if err := f.Close(); err != nil {
			log.Panicf("Failed to close file %s: %s\n", filename, err)
		}
	}()

	if err != nil {
		log.Fatalf("Failed to open file %s: %s\n", filename, err)
	}
	reader := csv.NewReader(f)
	// skip header
	reader.Read()

	for {

		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		t := NewTransactionFromValues(line)

		if strings.Contains(t.Description, "Transfer From") {
			log.Printf("skipping 'transfer from' transaction because firefly will handle the transfer records in both accounts")
			continue
		}

		ft := FireFlyTransaction{
			TransactionType: getType(t),
			Date:            t.Date,
			Amount:          t.Amount.Abs().String(),
			Category:        t.Category,
			DestinationName: t.DestinationAccount,
			SourceName:      t.SourceAccount,
			Tags:            t.Tags,
			Description:     t.Description,
		}

		// deposits require the dest account to be the source account for firefly
		if ft.TransactionType == "deposit" {
			ft.DestinationName = t.SourceAccount
			ft.SourceName = ""
		}

		payload := FireFlyTransactionRequest{
			ApplyRules:           true,
			ErrorIfDuplicateHash: !allowDuplicates,
			Transactions:         []FireFlyTransaction{ft},
		}

		err = client.AddTransaction(payload)

		if err != nil {
			log.Fatalf("Failed to add transaction to firefly: %s", err)
		}

	}

}
