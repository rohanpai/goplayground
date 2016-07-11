package main

import (
	&#34;encoding/json&#34;
	&#34;encoding/xml&#34;
	&#34;fmt&#34;
	&#34;github.com/PuerkitoBio/goquery&#34;
	&#34;regexp&#34;
	&#34;strings&#34;
)

const (
	ITEM        = &#34;div.b-parts-list-item.g-clrfix&#34;
	PRICE       = &#34;div.b-parts-list__price-block &gt; table.b-parts-list__price-block-table &gt; tbody &gt; tr &gt; td.b-parts-list__price-block-cell.b-parts-list__price-block-cell_price &gt; span.b-price.b-price_partslist&#34;
	NAME        = &#34;div.b-parts-list-item__name-block&#34;
	TITLE       = &#34;a.b-parts-list-item__name&#34;
	SKU         = &#34;div.b-parts-list-item__sku&#34;
	DESCRIPTION = &#34;div.b-parts-list-item__parameters &gt; span &gt; div&#34;
)

type Item struct {
	Name        string `xml:&#34;Name&#34; json:&#34;name&#34;`
	Price       string `xml:&#34;Price&#34; json:&#34;price&#34;`
	Link        string `xml:&#34;Link&#34; json:&#34;link&#34;`
	SKU         string `xml:&#34;SKU&#34; json:&#34;sku&#34;`
	Description string `xml:&#34;Description&#34;,json:&#34;description&#34;`
}

func Trim(text string) string {
	return regexp.MustCompile(&#34;[\t\r\n]&#34;).ReplaceAllString(text, &#34;&#34;)
}

func main() {
	var doc *goquery.Document
	var e error

	if doc, e = goquery.NewDocument(&#34;http://192.168.1.120:8000/mac.html&#34;); e != nil {
		panic(e.Error())
	}
	doc.Find(ITEM).Each(func(i int, s *goquery.Selection) {
		price := Trim(s.Find(PRICE).Text())
		name_block := s.Find(NAME)
		title_link := name_block.Find(TITLE)
		name := Trim(title_link.Text())
		link, _ := title_link.Attr(&#34;href&#34;)
		sku := strings.Replace(Trim(s.Find(SKU).Text()), &#34;SKU:&#34;, &#34;&#34;, -1)
		description := s.Find(DESCRIPTION).Text()
		item := &amp;Item{name, price, link, sku, description}
		item_json, _ := json.Marshal(item)
		item_xml, _ := xml.Marshal(item)
		fmt.Println(&#34;item:&#34;, item, string(item_json), string(item_xml))
	})
}
