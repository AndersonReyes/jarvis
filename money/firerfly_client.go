package money

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
)

const fireflyUrl = "http://192.168.1.212:8082/api"

var httpClient = &http.Client{}

type fireflyClient struct{}

func NewFireflyClient() fireflyClient {
	return fireflyClient{}
}

type FireFlyTransaction struct {
	TransactionType string   `json:"type"`
	Date            string   `json:"date"`
	Amount          string   `json:"amount"`
	Category        string   `json:"category_name"`
	DestinationName string   `json:"destination_name"`
	SourceName      string   `json:"source_name"`
	Tags            []string `json:"tags"`
	Description     string   `json:"description"`
}

type FireFlyTransactionRequest struct {
	ApplyRules           bool                 `json:"appy_rules"`
	ErrorIfDuplicateHash bool                 `json:"error_if_duplicate_hash"`
	Transactions         []FireFlyTransaction `json:"transactions"`
}

type FireFlyTransactioaErrorResponse struct {
	Message string `json:"message"`
	Errors  string `json:"errors"`
}

func (c fireflyClient) AddTransaction(r FireFlyTransactionRequest) error {
	payloadJson, err := json.Marshal(r)

	if err != nil {
		log.Printf("failed to convert %+v to json: %s\n", r, err)
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fireflyUrl+"/v1/transactions",
		bytes.NewBuffer(payloadJson))

	if err != nil {
		log.Printf("request building failed: %s\n", err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("FIREFLY_BEARER_TOKEN"))
	req.Header.Add("Accept", "application/json")

	resp, err := httpClient.Do(req)

	if err != nil {
		log.Printf("request to firefly failed: %s\n", err)
		return err
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	// var isDup = false
	// if resp.StatusCode == 422 || strings.Contains(string(body), "duplicate") {
	// 	isDup = true
	// }
	var buf bytes.Buffer
	err = json.Indent(&buf, body, "", "  ")

	if resp.StatusCode != 200 || err != nil {
		log.Printf("failed requesa: %+v\n", r)
		log.Printf("Error creating transaction with response Body:[%d] %s. Error: %s\n", resp.StatusCode, &buf, err)
		return errors.New("AddTransaction failed")
	}

	// log.Printf("API Response [%d]:\n%s\n", resp.StatusCode, &buf)
	// for _, t := range r.Transactions {
	// 	if isDup && r.ErrorIfDuplicateHash {
	// 		log.Printf("transaction exists. Skipping: %s\n", t)
	// 	} else {
	// 		log.Printf("created transaction: %s\n", t)
	// 	}
	// }

	return nil
}
