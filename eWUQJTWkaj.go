package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
)

const (
	ITEM        = "div.b-parts-list-item.g-clrfix"
	PRICE       = "div.b-parts-list__price-block > table.b-parts-list__price-block-table > tbody > tr > td.b-parts-list__price-block-cell.b-parts-list__price-block-cell_price > span.b-price.b-price_partslist"
	NAME        = "div.b-parts-list-item__name-block"
	TITLE       = "a.b-parts-list-item__name"
	SKU         = "div.b-parts-list-item__sku"
	DESCRIPTION = "div.b-parts-list-item__parameters > span > div"
)

type Item struct {
	Name        string `xml:"Name" json:"name"`
	Price       string `xml:"Price" json:"price"`
	Link        string `xml:"Link" json:"link"`
	SKU         string `xml:"SKU" json:"sku"`
	Description string `xml:"Description",json:"description"`
}

func Trim(text string) string {
	return regexp.MustCompile("[\t\r\n]").ReplaceAllString(text, "")
}

func main() {
	var doc *goquery.Document
	var e error

	if doc, e = goquery.NewDocument("http://192.168.1.120:8000/mac.html"); e != nil {
		panic(e.Error())
	}
	doc.Find(ITEM).Each(func(i int, s *goquery.Selection) {
		price := Trim(s.Find(PRICE).Text())
		name_block := s.Find(NAME)
		title_link := name_block.Find(TITLE)
		name := Trim(title_link.Text())
		link, _ := title_link.Attr("href")
		sku := strings.Replace(Trim(s.Find(SKU).Text()), "SKU:", "", -1)
		description := s.Find(DESCRIPTION).Text()
		item := &Item{name, price, link, sku, description}
		item_json, _ := json.Marshal(item)
		item_xml, _ := xml.Marshal(item)
		fmt.Println("item:", item, string(item_json), string(item_xml))
	})
}
