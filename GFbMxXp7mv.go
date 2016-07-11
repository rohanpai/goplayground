package main

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"strconv"
)

const (
	defHost   = "127.0.0.1"
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
	return s.Name != ""
}

type StrategyB struct {
	Hello string
	World string
}

func (s StrategyB) Print() string {
	return s.Hello
}

func (s StrategyB) IsValid() bool {
	return s.Hello != ""
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
	for _, impl := range []interface{}{&StrategyA{}, &StrategyB{}} {
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
	return fmt.Errorf("no suitable strategy type for %#v", m)
}

func main() {
	session, err := mgo.Dial(getHost() + ":" + strconv.Itoa(getPort()))
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("printer")
	err = c.DropCollection()
	if err != nil {
		panic(err)
	}
	str := []StrategyX{{StrategyA{"Ale", "+55 53 8116 9639"}}, {StrategyB{"Cla", "+55 53 8402 8510"}}}
	err = c.Insert(&Printer{str})
	if err != nil {
		panic(err)
	}

	var result Printer
	err = c.Find(nil).One(&result)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", result)
}
