package money

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"math"
	"os"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

var ReportCmd = &cobra.Command{
	Use:   "report [csv] [name]",
	Short: "generate html report of processed transaction csv",
	Long:  "generate html report of processed transaction csv",
	Args:  cobra.MinimumNArgs(2),
	Run:   generateReport,
}

type stats struct {
	In     decimal.Decimal `json:"in"`
	Out    decimal.Decimal `json:"out"`
	Net    decimal.Decimal `json:"net"`
	Budget decimal.Decimal `json:"budget"`
}

type JsonReport struct {
	DateStart        time.Time          `json:"date_start"`
	DateEndExclusive time.Time          `json:"date_end_exclusive"`
	Categories       map[Category]stats `json:"categories"`
	GlobalStats      stats              `json:"global_stats"`
	Transactions     []Transaction      `json:"transactions"`
}

func generateReport(cmd *cobra.Command, args []string) {
	log.SetPrefix("money.report: ")
	transactionsFileName := args[0]
	reportName := args[1]
	f, err := os.Open(transactionsFileName)
	defer func() {
		if err := f.Close(); err != nil {
			log.Panicf("Failed to close file %s: %s\n", transactionsFileName, err)
		}
	}()

	if err != nil {
		log.Fatalf("Failed to open file %s: %s\n", transactionsFileName, err)
	}

	reader := csv.NewReader(f)
	// skip reader
	reader.Read()

	var transactions []Transaction
	categoryStats := make(map[Category]stats)
	var globalStats stats
	maxDate := time.UnixMilli(0)
	minDate := time.UnixMilli(math.MaxInt64)

	zero := decimal.NewFromFloat(0.0)
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("failed to read next line: %s\n", err)
		}

		t := NewTransactionFromValues(line)

		entry, exists := categoryStats[t.Category]
		if !exists {
			entry = stats{}
		}

		entry.Net = entry.Net.Add(t.Amount)
		if t.Amount.GreaterThan(zero) {
			entry.In = entry.In.Add(t.Amount)
		} else {
			entry.Out = entry.Out.Add(t.Amount.Abs())
		}

		// set the budget if we have one set
		budget, ok := Budget[t.Category]
		if ok && entry.Budget.IsZero() {
			entry.Budget = budget
		}

		categoryStats[t.Category] = entry

		date, err := time.Parse(time.DateOnly, t.Date)
		if err != nil {
			log.Fatalf("failed to parse date %s: %s", t.Date, err)
		}

		if date.UnixMilli() < minDate.UnixMilli() {
			minDate = date
		}

		if date.UnixMilli() > maxDate.UnixMilli() {
			maxDate = date
		}

		transactions = append(transactions, t)
	}

	reportFile, err := os.OpenFile(reportName+".json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

	if err != nil {
		log.Fatalf("failed to create report.json: %s\n", err)
	}

	report := JsonReport{
		DateStart:        minDate,
		DateEndExclusive: maxDate.AddDate(0, 0, 1),
		Transactions:     transactions,
		Categories:       categoryStats,
		GlobalStats:      globalStats,
	}
	bytes, err := json.Marshal(report)

	if err != nil {
		log.Fatalf("failed to json encode stats: %s\n", err)
	}

	_, err = reportFile.Write(bytes)

	if err != nil {
		log.Fatalf("failed to write stats: %s\n", err)
	}

	// t := template.Must(template.New("report.tmpl").ParseFiles("templates/report.tmpl"))

	// htmlFile, err := os.OpenFile("report.html", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	// err = t.Execute(htmlFile, report)
	// if err != nil {
	// 	log.Fatalf("failed to write html report: %s \n", err)
	// }
}
