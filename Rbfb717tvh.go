package main

import (
       &#34;fmt&#34;
       &#34;time&#34;
       &#34;encoding/xml&#34;
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
        Amount int `xml:&#34;,innerxml&#34;`
}

type Plan struct {
        XMLName xml.Name `xml:&#34;plan&#34;`
        PlanCode string `xml:&#34;plan_code,omitempty&#34;`
        CreatedAt *time.Time `xml:&#34;created_at,omitempty&#34;`
        UnitAmountInCents CurrencyArray `xml:&#34;unit_amount_in_cents&#34;`
        SetupFeeInCents CurrencyArray `xml:&#34;setup_in_cents&#34;`
}

func main() {
	fmt.Println(&#34;Hello, playground&#34;)
	p := Plan{PlanCode:&#34;test&#34;}
	p.UnitAmountInCents.AddCurrency(&#34;USD&#34;,4000)
        if xmlstring, err := xml.MarshalIndent(p, &#34;&#34;, &#34;    &#34;); err == nil {
                xmlstring = []byte(xml.Header &#43; string(xmlstring))
                fmt.Printf(&#34;%s\n&#34;,xmlstring)
        }
}
