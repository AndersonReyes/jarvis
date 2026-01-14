package money

import (
	"encoding/csv"
	"io"
	"log"
	"os"

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
)

var client = NewFireflyClient()

func init() {
	FireflyCmd.Flags().BoolVarP(&importTransactionsFlag, "transactions", "t",
		false, "import transactions from csv")
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
	if t.Amount.LessThan(decimal.NewFromFloat(0.0)) {
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
		ft := FireFlyTransaction{
			TransactionType: getType(t),
			Date:            t.Date,
			Amount:          t.Amount.String(),
			Category:        t.Category,
			DestinationName: t.DestinationAccount,
			SourceName:      t.SourceAccount,
			Tags:            t.Tags,
			Description:     t.Description,
		}

		if ft.DestinationName == "" {
			ft.DestinationName = "noname"
		}

		if ft.SourceName == "" {
			ft.SourceName = "noname"
		}

		payload := FireFlyTransactionRequest{
			ApplyRules:           true,
			ErrorIfDuplicateHash: true,
			Transactions:         []FireFlyTransaction{ft},
		}

		err = client.AddTransaction(payload)

		if err != nil {
			log.Fatalf("Failed to add transaction to firefly: %s", err)
		}

	}

}
