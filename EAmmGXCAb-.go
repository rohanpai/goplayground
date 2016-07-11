package main
 
import (
	&#34;fmt&#34;
	&#34;sort&#34;
	&#34;math/rand&#34;
	&#34;time&#34;
	&#34;strings&#34;
)
 
// Kart
type Card struct {
	Suite      string
	Rank       string
}
 
// Kartı sembol olarak ver
func (c *Card) Symbol() string {
	return c.Suite[:1] &#43; c.Rank
}
 
func (c *Card) RankValue() int {
	return sort.SearchStrings(RANKS, c.Rank)
}
 
func (c *Card) SuiteValue() int {
	return sort.SearchStrings(SUITES, c.Suite)
}
 
// Kartı sembolle göster
func (c *Card) ToString() string {
	return GLYPH[c.Symbol()]
}
 
// Deste
type Deck struct {
	Cards []Card
}
 
// Desteden el çıkar
func (d *Deck) Hand() Hand {
	// pop
	hand_cards := d.Cards[0:DEFAULT_HAND_SIZE]
	// pop sonrası kalan deste
	d.Cards = d.Cards[DEFAULT_HAND_SIZE:]
	return Hand{hand_cards}
}
 
// Desteyi karıştır
func (d *Deck) Shuffle() *Deck {
	src := d.Cards
	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(len(d.Cards))
	for i, v := range perm {
		d.Cards[v] = src[i]
	}
	return d
}
 
// El
type Hand struct {
	Cards []Card
}
 
func (h *Hand) Deal() Card {
	// pop işlemi
	card := h.Cards[len(h.Cards)-1:]
	h.Cards = h.Cards[0:len(h.Cards)-1]
	return card[0]
}
 
// Eli sembollerle göster
func (h *Hand) ToString() string {
	var cardsSymbol []string
	for _, card := range h.Cards {
		cardsSymbol = append(cardsSymbol, GLYPH[card.Symbol()])
	}
	return strings.Join(cardsSymbol, &#34; &#34;)
}
 
var (
	SUITES []string = []string{&#34;club&#34;, &#34;diamondi&#34;, &#34;heart&#34;, &#34;spade&#34;}
	RANKS []string = []string{&#34;1&#34;, &#34;2&#34;, &#34;3&#34;, &#34;4&#34;, &#34;5&#34;, &#34;6&#34;, &#34;7&#34;, &#34;8&#34;, &#34;9&#34;, &#34;10&#34;, &#34;j&#34;, &#34;q&#34;, &#34;k&#34;}
	DEFAULT_HAND_SIZE = 4 // Öntanımlı el kaç kart
	GLYPH map[string]string = map[string]string {
		&#34;s1&#34;:  &#34;\U0001f0a1&#34;, &#34;h1&#34;:  &#34;\U0001f0b1&#34;, &#34;d1&#34;:  &#34;\U0001f0c1&#34;, &#34;c1&#34;:  &#34;\U0001f0d1&#34;,
		&#34;s2&#34;:  &#34;\U0001f0a2&#34;, &#34;h2&#34;:  &#34;\U0001f0b2&#34;, &#34;d2&#34;:  &#34;\U0001f0c2&#34;, &#34;c2&#34;:  &#34;\U0001f0d2&#34;,
		&#34;s3&#34;:  &#34;\U0001f0a3&#34;, &#34;h3&#34;:  &#34;\U0001f0b3&#34;, &#34;d3&#34;:  &#34;\U0001f0c3&#34;, &#34;c3&#34;:  &#34;\U0001f0d3&#34;,
		&#34;s4&#34;:  &#34;\U0001f0a4&#34;, &#34;h4&#34;:  &#34;\U0001f0b4&#34;, &#34;d4&#34;:  &#34;\U0001f0c4&#34;, &#34;c4&#34;:  &#34;\U0001f0d4&#34;,
		&#34;s5&#34;:  &#34;\U0001f0a5&#34;, &#34;h5&#34;:  &#34;\U0001f0b5&#34;, &#34;d5&#34;:  &#34;\U0001f0c5&#34;, &#34;c5&#34;:  &#34;\U0001f0d5&#34;,
		&#34;s6&#34;:  &#34;\U0001f0a6&#34;, &#34;h6&#34;:  &#34;\U0001f0b6&#34;, &#34;d6&#34;:  &#34;\U0001f0c6&#34;, &#34;c6&#34;:  &#34;\U0001f0d6&#34;,
		&#34;s7&#34;:  &#34;\U0001f0a7&#34;, &#34;h7&#34;:  &#34;\U0001f0b7&#34;, &#34;d7&#34;:  &#34;\U0001f0c7&#34;, &#34;c7&#34;:  &#34;\U0001f0d7&#34;,
		&#34;s8&#34;:  &#34;\U0001f0a8&#34;, &#34;h8&#34;:  &#34;\U0001f0b8&#34;, &#34;d8&#34;:  &#34;\U0001f0c8&#34;, &#34;c8&#34;:  &#34;\U0001f0d8&#34;,
		&#34;s9&#34;:  &#34;\U0001f0a9&#34;, &#34;h9&#34;:  &#34;\U0001f0b9&#34;, &#34;d9&#34;:  &#34;\U0001f0c9&#34;, &#34;c9&#34;:  &#34;\U0001f0d9&#34;,
		&#34;s10&#34;: &#34;\U0001f0aa&#34;, &#34;h10&#34;: &#34;\U0001f0ba&#34;, &#34;d10&#34;: &#34;\U0001f0ca&#34;, &#34;c10&#34;: &#34;\U0001f0da&#34;,
		&#34;sj&#34;:  &#34;\U0001f0ab&#34;, &#34;hj&#34;:  &#34;\U0001f0bb&#34;, &#34;dj&#34;:  &#34;\U0001f0cb&#34;, &#34;cj&#34;:  &#34;\U0001f0db&#34;,
		&#34;sq&#34;:  &#34;\U0001f0ad&#34;, &#34;hq&#34;:  &#34;\U0001f0bd&#34;, &#34;dq&#34;:  &#34;\U0001f0cd&#34;, &#34;cq&#34;:  &#34;\U0001f0dd&#34;,
		&#34;sk&#34;:  &#34;\U0001f0ae&#34;, &#34;hk&#34;:  &#34;\U0001f0be&#34;, &#34;dk&#34;:  &#34;\U0001f0ce&#34;, &#34;ck&#34;:  &#34;\U0001f0de&#34;,
	}
)
 
 
func main() {
	// Kartları şimdi üret
	var CARDS []Card
	for _, suite := range SUITES {
		for _, rank := range RANKS {
			card := Card{suite, rank}
			CARDS = append(CARDS, card)
		}
	}
 
	// Yeni deste
	d := Deck{CARDS}
	// Desteyi karıştır
	d.Shuffle()
 
	for i := 0; i &lt; (52 / DEFAULT_HAND_SIZE); i&#43;&#43; {
		h := d.Hand() // Desteden eller
		fmt.Println(h.ToString())
		for j := 0; j &lt; DEFAULT_HAND_SIZE; j&#43;&#43; {
			card := h.Deal() // Ellerden tek tek kart
			fmt.Println(card.ToString())
		}
	}
}