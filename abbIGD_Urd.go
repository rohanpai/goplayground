package main

import (
	//&#34;encoding/csv&#34;
	&#34;github.com/davecgh/go-spew/spew&#34;
	&#34;log&#34;
	&#34;net/http&#34;
)

func main() {

	req, err := http.NewRequest(&#34;GET&#34;, &#34;http://rbl-check.org/rbl_api.php?ipaddress=95.89.93.164&#34;, nil)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	if resp.StatusCode != 200 {
		log.Fatalln(&#34;Status:&#34;, resp.Status)
	}
	defer resp.Body.Close()

	spew.Dump(resp)

}
