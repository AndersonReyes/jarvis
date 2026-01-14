package money

import (
	"encoding/csv"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

var MoneyCmd = &cobra.Command{
	Use:   "money [dir]",
	Short: "manage yo money, all csvs in [dir]",
	Long:  "clean, process, generate financial data for tracking expenses. Note: uses the filename as account name",
	Args:  cobra.MinimumNArgs(1),
	Run:   run,
}

var numberValueRegex = regexp.MustCompile("[$,]")

func processDiscoveryAccounts(account string, row []string) (*Transaction, error) {

	var amountStr string

	if strings.ToLower(row[2]) == "debit" {
		amountStr = "-" + row[3]
	} else {
		amountStr = row[4]
	}

	amount, err := decimal.NewFromFormattedString(amountStr, numberValueRegex)
	if err != nil {
		log.Fatalf("Failed to parse debit/credit for %v: %v", row, err)
		return nil, err
	}

	balance, err := decimal.NewFromFormattedString(row[5], numberValueRegex)
	if err != nil {
		log.Fatalf("Failed to parse balance for %v: %v", row, err)
		return nil, err
	}

	dateParts := strings.Split(row[0], "/")
	cleanDate := strings.Join([]string{dateParts[2], dateParts[0], dateParts[1]}, "-")
	externalAcc := AccountIdToFullName[account]
	tags := []string{" #bank-discover-" + account}

	record := &Transaction{
		Date:          cleanDate,
		Description:   row[1],
		Amount:        amount,
		Balance:       balance,
		Bank:          "discover",
		SourceAccount: externalAcc,
		Tags:          tags,
		Category:      CategoryUnknown,
	}

	SetCategory("", record)
	SetPayeeAndDestinationAccount(record)

	return record, nil

}

func processCapitalOneAccounts(row []string) (*Transaction, error) {

	var amountStr string

	if row[len(row)-2] != "" {
		amountStr = "-" + row[len(row)-2]
	} else {
		amountStr = row[len(row)-1]
	}

	amount, err := decimal.NewFromFormattedString(amountStr, numberValueRegex)
	if err != nil {
		log.Fatalf("Failed to parse debit/credit for %v: %v", row, err)
		return nil, err
	}

	externalAcc := AccountIdToFullName["8444"]
	tags := []string{" #bank-capitalone-8444"}
	record := &Transaction{
		Date:          row[0],
		Description:   row[3],
		Amount:        amount,
		Tags:          tags,
		Balance:       decimal.NewFromFloat32(0.0),
		Bank:          "capitalone",
		SourceAccount: externalAcc,
		Category:      "N/A",
	}
	SetCategory(row[4], record)
	SetPayeeAndDestinationAccount(record)

	return record, nil

}

const output = "transactions.csv"

func run(cmd *cobra.Command, args []string) {
	log.SetPrefix("money")
	for _, dir := range args {
		var files []string

		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				log.Fatalf("failed to walk directory for files: %v", err)
				return err
			}

			if !d.IsDir() && strings.HasSuffix(d.Name(), ".csv") {
				files = append(files, filepath.Join(dir, d.Name()))
			}

			return nil
		})

		if err != nil {
			log.Panicf("failed to iterate directory %+v: %s", args, err)

		}

		outFile, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		defer func() {
			if err := outFile.Close(); err != nil {
				log.Panicf("Failed to close file %s: %s\n", output, err)
			}
		}()
		if err != nil {
			log.Fatalf("failed to open %s for writing: %s", output, err)
		}

		writer := csv.NewWriter(outFile)
		err = writer.Write(TransactionColumns[:])
		if err != nil {
			log.Fatalf("Failed to write header to %s: %s", output, err)
		}

		for _, filename := range files {

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

			log.Printf("Processing %s now\n", filename)
			for {

				line, err := reader.Read()
				if err == io.EOF {
					break
				}

				log.Printf("Processing line %s", line)

				var record *Transaction

				if strings.Contains(filename, "capitalone") {
					record, err = processCapitalOneAccounts(line)
				} else {
					// uses the filename as account name
					account := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
					record, err = processDiscoveryAccounts(account, line)
				}

				if err != nil {
					log.Fatalf("Failed to create record: %s", err)
				}
				// log.Printf("%+v\n", record)

				err = writer.Write([]string{
					record.Date, record.SourceAccount, record.Category,
					record.Amount.String(), record.Description, record.Payee,
					record.DestinationAccount,
				})
				if err != nil {
					log.Fatalf("error writing record %+v to csv: %s", record, err)
				}
			}
		}

		writer.Flush()
	}
}
