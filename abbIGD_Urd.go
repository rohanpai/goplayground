package main

import (
	//"encoding/csv"
	"github.com/davecgh/go-spew/spew"
	"log"
	"net/http"
)

func main() {

	req, err := http.NewRequest("GET", "http://rbl-check.org/rbl_api.php?ipaddress=95.89.93.164", nil)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	if resp.StatusCode != 200 {
		log.Fatalln("Status:", resp.Status)
	}
	defer resp.Body.Close()

	spew.Dump(resp)

}
