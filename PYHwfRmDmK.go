package main

import (
	&#34;errors&#34;
	&#34;fmt&#34;
	&#34;time&#34;
)

type Message struct {
	Data      []byte
	MimeType  string
	Timestamp time.Time
}

type TwitterSource struct {
	Username string
}

type SkypeSource struct {
	Login         string
	MSBackdoorKey string
}

// Finder represents any source for words lookup.
type Finder interface {
	Find(word string) ([]Message, error)
}

func (s TwitterSource) Find(word string) ([]Message, error) {
	return s.searchAPICall(s.Username, word)
}

func (s SkypeSource) Find(word string) ([]Message, error) {
	return s.searchSkypeServers(s.MSBackdoorKey, s.Login, word)
}

type Sources []Finder

func (s Sources) SearchWords(word string) []Message {
	var messages []Message
	for _, source := range s {
		msgs, err := source.Find(word)
		if err != nil {
			fmt.Println(&#34;WARNING:&#34;, err)
			continue
		}
		messages = append(messages, msgs...)
	}

	return messages
}

var (
	sources = Sources{
		TwitterSource{
			Username: &#34;@rickhickey&#34;,
		},
		SkypeSource{
			Login:         &#34;rich.hickey&#34;,
			MSBackdoorKey: &#34;12345&#34;,
		},
	}

	person = Person{
		FullName: &#34;Rick Hickey&#34;,
		Sources:  sources,
	}
)

type Person struct {
	FullName string
	Sources
}

func main() {
	msgs := person.SearchWords(&#34;если бы бабушка&#34;)
	fmt.Println(msgs)
}

func (s TwitterSource) searchAPICall(username, word string) ([]Message, error) {
	return []Message{
		Message{
			Data:      ([]byte)(&#34;Remember, remember, the fifth of November, если бы бабушка...&#34;),
			MimeType:  &#34;text/plain&#34;,
			Timestamp: time.Now(),
		},
	}, nil
}

func (s SkypeSource) searchSkypeServers(key, login, word string) ([]Message, error) {
	return []Message{}, errors.New(&#34;NSA can&#39;t read your skype messages ;)&#34;)
}

func (m Message) String() string {
	return string(m.Data) &#43; &#34; @ &#34; &#43; m.Timestamp.Format(time.RFC822)
}