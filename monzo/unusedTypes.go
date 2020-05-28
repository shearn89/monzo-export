package monzo

import "encoding/xml"

type LedgerBalance struct {
	XMLName xml.Name `xml:"LEDGERBAL"`
	Amount float32 `xml:"BALAMT"`
	Date int `xml:"DTASOF"`
}
type AvailableBalance struct {
	XMLName xml.Name `xml:"AVAILBAL"`
	Amount float32 `xml:"BALAMT"`
	Date int `xml:"DTASOF"`
}


type Status struct {
	XMLName xml.Name `xml:"STATUS"`
	Code int `xml:"CODE"`
	Severity string `xml:"SEVERITY"`
}
