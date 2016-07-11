package main

import (
	"encoding/json"
	"fmt"
)

type request struct {
	Operations map[string]Op `json:"operations"`
}
type Op struct {
	Opp  Operation `json:"opp"`
	Test string    `json:"test"`
}
type Operation struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func (o *Operation) UnmarshalJSON(b []byte) error {
	type xoperation Operation
	xo := &xoperation{Width: 500, Height: 50}
	if err := json.Unmarshal(b, xo); err != nil {
		return err
	}
	*o = Operation(*xo)
	return nil
}

func main() {
	jsonStr := `{
			"operations": {
				"001": {
			 		"test":"test1",
			 		"opp":{
					"width": 101
			 		}
				},
				"002": {
			 		"test":"test2",
			 		"opp":{
					"width": 102
			 		}
				}
			}
		}`
	req := request{}
	json.Unmarshal([]byte(jsonStr), &req)
	fmt.Println(req)
}
