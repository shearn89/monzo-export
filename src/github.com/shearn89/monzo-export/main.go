package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"io/ioutil"
	"encoding/json"
	"encoding/xml"
	"net/http"

	monzo "github.com/shearn89/monzo-export/monzo"
)

func getTransactions(output *monzo.MonzoTransactions, string lastTransId) error {
	// since := "2017-06-16T00:00:00Z"
	token := os.Getenv("MONZO_ACCESS_TOKEN")
	accId := os.Getenv("MONZO_ACCOUNT_ID")

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.monzo.com/transactions", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	q := req.URL.Query()
	q.Add("account_id", accId)
	q.Add("expand[]", "merchant")
	if lastTransId != nil {
		q.Add("since", lastTransId)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, output)
}

func getAccounts() (string, error) {
	token := os.Getenv("MONZO_ACCESS_TOKEN")

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.monzo.com/accounts", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var output monzo.MonzoAccounts
	err = json.Unmarshal(data, &output)
	if err != nil {
		return "", err
	}
	return output.Accounts[0].ID, nil
}

func main() {
	lastTransId := nil
	dat, err := ioutil.ReadAll(".last_export")
	if err != nil {
		fmt.Println("Unable to read .last_export file, getting all transactions")
	} else {
		lastTransId = string(dat)
	}

	var trans monzo.MonzoTransactions
	getTransactions(&trans, lastTransId)
	fmt.Printf("There are %d transactions to process.\n", len(trans.Transactions))
	if len(trans.Transactions) == 0 {
		fmt.Println("Exiting.")
		os.Exit(0)
	}
	
	dateStart,err := strconv.Atoi(trans.Transactions[0].Created.Format("20060102"))
	if err != nil {
		log.Fatal(err)
	}
	maxId := len(trans.Transactions)-1
	dateEnd, err := strconv.Atoi(trans.Transactions[maxId].Created.Format("20060102"))
	if err != nil {
		log.Fatal(err)
	}
	lastTransId = trans.Transactions[maxId].ID
	
	xmlTrans := []monzo.Transaction{}
	for _,t := range trans.Transactions {
		var transType string
		if t.IsLoad {
			transType = "CREDIT"
		} else {
			transType = "POS"
		}
		date, err := strconv.Atoi(t.Created.Format("20060102"))
		if err != nil {
			log.Fatal(err)
		}
		amount := float32(t.Amount)/100
		id := t.ID
		name := t.Merchant.Name
		note := strings.Split(t.Description, "  ")[0]

		outTrans := monzo.Transaction{
			Type: transType,
			Date: date,
			Amount: amount,
			Id: id,
			Name: name,
			Note: note,
		}
		xmlTrans = append(xmlTrans, outTrans)
	}

	acctNumber, err := getAccounts()
	if err != nil {
		log.Fatal(err)
	}
	
	xmlOfx := monzo.OFX{
		BankMessages: monzo.BankMessages{
			StatementRecord: monzo.StatementRecord{
				StatementTransactions: monzo.StatementTransactions{
					BankAccount: monzo.BankAccount{
						AccountNumber: acctNumber,
						AccountType: "CHECKING",
					},
					TransactionList: monzo.TransactionList{
						Start: dateStart,
						End: dateEnd,
						Transactions: xmlTrans,
					},
				},
			},
		},
	}

	data, err := xml.MarshalIndent(xmlOfx, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	xmlFileName := fmt.Sprintf("monzo-export-%d.ofx", int32(time.Now().Unix()))
	err = ioutil.WriteFile(xmlFileName, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
	
	err = ioutil.WriteFile(".last_export", []byte(lastTransId), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

