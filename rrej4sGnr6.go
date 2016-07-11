package main

import (
	&#34;encoding/json&#34;
	&#34;fmt&#34;
)

type request struct {
	Operations map[string]Op `json:&#34;operations&#34;`
}
type Op struct {
	Opp  Operation `json:&#34;opp&#34;`
	Test string    `json:&#34;test&#34;`
}
type Operation struct {
	Width  int `json:&#34;width&#34;`
	Height int `json:&#34;height&#34;`
}

func (o *Operation) UnmarshalJSON(b []byte) error {
	type xoperation Operation
	xo := &amp;xoperation{Width: 500, Height: 50}
	if err := json.Unmarshal(b, xo); err != nil {
		return err
	}
	*o = Operation(*xo)
	return nil
}

func main() {
	jsonStr := `{
			&#34;operations&#34;: {
				&#34;001&#34;: {
			 		&#34;test&#34;:&#34;test1&#34;,
			 		&#34;opp&#34;:{
					&#34;width&#34;: 101
			 		}
				},
				&#34;002&#34;: {
			 		&#34;test&#34;:&#34;test2&#34;,
			 		&#34;opp&#34;:{
					&#34;width&#34;: 102
			 		}
				}
			}
		}`
	req := request{}
	json.Unmarshal([]byte(jsonStr), &amp;req)
	fmt.Println(req)
}
