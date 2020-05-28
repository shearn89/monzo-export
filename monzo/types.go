package monzo

import (
	"time"
	"encoding/xml"
)

type MonzoTransactions struct {
	Transactions []MonzoTransaction
}

type MonzoAccounts struct {
	Accounts []struct {
		ID string `json:"id"`
		Created time.Time `json:"created"`
		Description string `json:"description"`
		Type string `json:"type"`
	} `json:"accounts"`
}

type MonzoTransaction struct {
	ID string `json:"id"`
	Created time.Time `json:"created"`
	Description string `json:"description"`
	Amount int `json:"amount"`
	Currency string `json:"currency"`
	Merchant struct {
		ID string `json:"id"`
		GroupID string `json:"group_id"`
		Created time.Time `json:"created"`
		Name string `json:"name"`
		Logo string `json:"logo"`
		Emoji string `json:"emoji"`
		Category string `json:"category"`
		Online bool `json:"online"`
		Atm bool `json:"atm"`
		Address struct {
			ShortFormatted string `json:"short_formatted"`
			Formatted string `json:"formatted"`
			Address string `json:"address"`
			City string `json:"city"`
			Region string `json:"region"`
			Country string `json:"country"`
			Postcode string `json:"postcode"`
			Latitude float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			ZoomLevel int `json:"zoom_level"`
			Approximate bool `json:"approximate"`
		} `json:"address"`
		Updated time.Time `json:"updated"`
		Metadata struct {
			CreatedForMerchant string `json:"created_for_merchant"`
			CreatedForTransaction string `json:"created_for_transaction"`
			EnrichedFromSettlement string `json:"enriched_from_settlement"`
			FoursquareCategory string `json:"foursquare_category"`
			FoursquareCategoryIcon string `json:"foursquare_category_icon"`
			FoursquareID string `json:"foursquare_id"`
			FoursquareWebsite string `json:"foursquare_website"`
			GooglePlacesIcon string `json:"google_places_icon"`
			GooglePlacesID string `json:"google_places_id"`
			GooglePlacesName string `json:"google_places_name"`
			SuggestedName string `json:"suggested_name"`
			SuggestedTags string `json:"suggested_tags"`
			TwitterID string `json:"twitter_id"`
			Website string `json:"website"`
		} `json:"metadata"`
		DisableFeedback bool `json:"disable_feedback"`
	} `json:"merchant"`
	Notes string `json:"notes"`
	Metadata struct {
	} `json:"metadata"`
	AccountBalance int `json:"account_balance"`
	Attachments []interface{} `json:"attachments"`
	Category string `json:"category"`
	IsLoad bool `json:"is_load"`
	Settled time.Time `json:"settled"`
	LocalAmount int `json:"local_amount"`
	LocalCurrency string `json:"local_currency"`
	Updated time.Time `json:"updated"`
	AccountID string `json:"account_id"`
	Counterparty struct {
	} `json:"counterparty"`
	Scheme string `json:"scheme"`
	DedupeID string `json:"dedupe_id"`
	Originator bool `json:"originator"`
	IncludeInSpending bool `json:"include_in_spending"`
}

type Transaction struct {
	XMLName xml.Name `xml:"STMTTRN"`
	Type string `xml:"TRNTYPE"`
	Date int `xml:"DTPOSTED"`
	Amount float32 `xml:"TRNAMT"`
	Id string `xml:"FITID"`
	Name string `xml:"NAME"`
	Note string `xml:"MEMO"`
}

type BankAccount struct {
	XMLName xml.Name `xml:"BANKACCTFROM"`
	// SortCode int `xml:"BANKID"`
	AccountNumber string `xml:"ACCTID"`
	AccountType string `xml:"ACCTTYPE"`
}
type TransactionList struct {
	XMLName xml.Name `xml:"BANKTRANLIST"`
	Start int `xml:"DTSTART"`
	End int `xml:"DTEND"`
	Transactions []Transaction
}
type StatementTransactions struct {
	XMLName xml.Name `xml:"STMTRS"`
	// Currency string `xml:"CURDEF"`
	BankAccount BankAccount
	TransactionList TransactionList
	// LedgerBalance
	// AvailableBalance
}

type StatementRecord struct {
	XMLName xml.Name `xml:"STMTTRNRS"`
	// TrnId int `xml:"TRNUID"`
	// Status
	StatementTransactions StatementTransactions
}

type BankMessages struct {
	XMLName xml.Name `xml:"BANKMSGSRSV1"`
	StatementRecord StatementRecord
}

type OFX struct {
	XMLName xml.Name `xml:"OFX"`
	BankMessages BankMessages
}

