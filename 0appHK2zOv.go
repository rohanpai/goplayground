package main

import (
	&#34;encoding/xml&#34;
	&#34;fmt&#34;
	&#34;strings&#34;
)

type Reading struct {
	RType		string	`xml:&#34;r_type,attr&#34;`
	ReadingString	string	`xml:&#34;,innerxml&#34;`
}

type Meaning struct {
	MLang		string	`xml:&#34;m_lang,attr&#34;`
	MeaningString	string	`xml:&#34;,innerxml&#34;`
}

type DicRef struct {
	DrType		string	`xml:&#34;dr_type,attr&#34;`
	RefNumber	string	`xml:&#34;,innerxml&#34;`
}

type Kanji struct {
	Literal		string		`xml:&#34;literal&#34;`
	Grade		int		`xml:&#34;misc&gt;grade&#34;`
	StrokeCount	int		`xml:&#34;misc&gt;stroke_count&#34;`
	Freq		int		`xml:&#34;misc&gt;freq&#34;`
	JLPT		int		`xml:&#34;misc&gt;jlpt&#34;`
	DicRefs		[]DicRef	`xml:&#34;dic_number&gt;dic_ref&#34;`
	Readings	[]Reading	`xml:&#34;reading_meaning&gt;rmgroup&gt;reading&#34;`
	Meanings	[]Meaning	`xml:&#34;reading_meaning&gt;rmgroup&gt;meaning&#34;`
	Nanori		[]string	`xml:&#34;nanori&#34;`
}

var theXML string = `&lt;?xml version=&#34;1.0&#34; encoding=&#34;UTF-8&#34;?&gt;
&lt;!DOCTYPE kanjidic2 [
&lt;!ELEMENT variant (#PCDATA)&gt;
&lt;!ATTLIST variant var_type CDATA #REQUIRED&gt;
	&lt;!-- 
	oneill - Japanese Names (O&#39;Neill) - numeric
	--&gt;

]&gt;
&lt;header&gt;

&lt;file_version&gt;4&lt;/file_version&gt;
&lt;database_version&gt;2011-536&lt;/database_version&gt;
&lt;date_of_creation&gt;2012-06-19&lt;/date_of_creation&gt;
&lt;/header&gt;
&lt;character&gt;
&lt;literal&gt;本&lt;/literal&gt;
&lt;misc&gt;
&lt;grade&gt;1&lt;/grade&gt;
&lt;stroke_count&gt;5&lt;/stroke_count&gt;
&lt;variant var_type=&#34;jis208&#34;&gt;52-81&lt;/variant&gt;
&lt;freq&gt;10&lt;/freq&gt;
&lt;jlpt&gt;4&lt;/jlpt&gt;
&lt;/misc&gt;
&lt;dic_number&gt;
&lt;dic_ref dr_type=&#34;nelson_c&#34;&gt;96&lt;/dic_ref&gt;
&lt;/dic_number&gt;
&lt;query_code&gt;
&lt;q_code qc_type=&#34;skip&#34;&gt;4-5-3&lt;/q_code&gt;
&lt;/query_code&gt;
&lt;reading_meaning&gt;
&lt;rmgroup&gt;
&lt;reading r_type=&#34;ja_on&#34;&gt;ホン&lt;/reading&gt;
&lt;reading r_type=&#34;ja_kun&#34;&gt;もと&lt;/reading&gt;
&lt;meaning&gt;book&lt;/meaning&gt;
&lt;meaning m_lang=&#34;fr&#34;&gt;livre&lt;/meaning&gt;
&lt;/rmgroup&gt;
&lt;nanori&gt;まと&lt;/nanori&gt;
&lt;/reading_meaning&gt;
&lt;/character&gt;
&lt;/kanjidic2&gt;
`

func ParseKanjiDic2() (kanjiList []Kanji) {
	decoder := xml.NewDecoder(strings.NewReader(theXML))
	for {
		token, _ := decoder.Token()
		//fmt.Println(token)
		if token == nil {
			break
		}
		switch startElement := token.(type) {
		case xml.StartElement:
			if startElement.Name.Local == &#34;character&#34; {
				var kanji Kanji
				decoder.DecodeElement(&amp;kanji, &amp;startElement)
				kanjiList = append(kanjiList, kanji)
			}
		}
	}
	return
}

func main() {
	kanjidic := ParseKanjiDic2()
	fmt.Println(kanjidic)
}
