package main

import (
       "fmt"
       "time"
       "encoding/xml"
)

type CurrencyArray struct {
        CurrencyList []Currency
}

func (c *CurrencyArray) AddCurrency(currency string, amount int) {
        newc := Currency{Amount:amount}
        newc.XMLName.Local = currency
        c.CurrencyList = append(c.CurrencyList, newc)
}

type Currency struct {
        XMLName xml.Name 
        Amount int `xml:",innerxml"`
}

type Plan struct {
        XMLName xml.Name `xml:"plan"`
        PlanCode string `xml:"plan_code,omitempty"`
        CreatedAt *time.Time `xml:"created_at,omitempty"`
        UnitAmountInCents CurrencyArray `xml:"unit_amount_in_cents"`
        SetupFeeInCents CurrencyArray `xml:"setup_in_cents"`
}

func main() {
	fmt.Println("Hello, playground")
	p := Plan{PlanCode:"test"}
	p.UnitAmountInCents.AddCurrency("USD",4000)
        if xmlstring, err := xml.MarshalIndent(p, "", "    "); err == nil {
                xmlstring = []byte(xml.Header + string(xmlstring))
                fmt.Printf("%s\n",xmlstring)
        }
}
