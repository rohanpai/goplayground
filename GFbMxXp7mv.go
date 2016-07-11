package main

import (
	&#34;fmt&#34;
	&#34;labix.org/v2/mgo&#34;
	&#34;labix.org/v2/mgo/bson&#34;
	&#34;strconv&#34;
)

const (
	defHost   = &#34;127.0.0.1&#34;
	defPort   = 27017
)

func getHost() string {
	return defHost
}

func getPort() int {
	return defPort
}

type Strategy interface {
	IsValid() bool
	Print() string
}

type StrategyA struct {
	Name  string
	Phone string
}

func (s StrategyA) Print() string {
	return s.Name
}

func (s StrategyA) IsValid() bool {
	return s.Name != &#34;&#34;
}

type StrategyB struct {
	Hello string
	World string
}

func (s StrategyB) Print() string {
	return s.Hello
}

func (s StrategyB) IsValid() bool {
	return s.Hello != &#34;&#34;
}

type Printer struct {
	Strategy []StrategyX
}

type StrategyX struct {
	Strategy
}

func (s StrategyX) GetBSON() (interface{}, error) {
	return s.Strategy, nil
}

func (s *StrategyX) SetBSON(raw bson.Raw) error {
	for _, impl := range []interface{}{&amp;StrategyA{}, &amp;StrategyB{}} {
		err := raw.Unmarshal(impl)
		if err != nil {
			return err
		}
		if impl.(Strategy).IsValid() {
			s.Strategy = impl.(Strategy)
			return nil
		}
	}
	m := bson.M{}
	err := raw.Unmarshal(m)
	if err != nil {
		return err
	}
	return fmt.Errorf(&#34;no suitable strategy type for %#v&#34;, m)
}

func main() {
	session, err := mgo.Dial(getHost() &#43; &#34;:&#34; &#43; strconv.Itoa(getPort()))
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(&#34;test&#34;).C(&#34;printer&#34;)
	err = c.DropCollection()
	if err != nil {
		panic(err)
	}
	str := []StrategyX{{StrategyA{&#34;Ale&#34;, &#34;&#43;55 53 8116 9639&#34;}}, {StrategyB{&#34;Cla&#34;, &#34;&#43;55 53 8402 8510&#34;}}}
	err = c.Insert(&amp;Printer{str})
	if err != nil {
		panic(err)
	}

	var result Printer
	err = c.Find(nil).One(&amp;result)
	if err != nil {
		panic(err)
	}

	fmt.Printf(&#34;%v&#34;, result)
}
