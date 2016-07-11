package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	BASEURL string = "http://www.stanford.edu/class/cs193c"
)

func main() {
	resp, err1 := http.Get(BASEURL + "/lectures.html")
	panicIf(err1)

	defer resp.Body.Close()

	body, err2 := ioutil.ReadAll(resp.Body)
	panicIf(err2)

	bodyString := string(body)

	linkList := make([]string, 0)
	getLinks(bodyString, &linkList)

	fmt.Println("Fetching Urls...")

	for _, link := range linkList {
		fmt.Println(link)
	}

	fmt.Println("Done.")

	os.Exit(0)
}

func getLinks(webpage string, linkList *[]string) {
	gotLink, i, err := getLink(webpage)
	if err != nil {
		return
	}

	if strings.HasSuffix(gotLink, ".zip") {
		*linkList = append(*linkList, gotLink)
	}

	if strings.HasSuffix(gotLink, ".pdf") {
		*linkList = append(*linkList, gotLink)
	}

	getLinks(webpage[i:], linkList)
}

func getLink(webpage string) (string, int, error) {
	target1 := strings.Index(webpage, "a href=\"")
	if target1 == -1 {
		return "", -1, errors.New("No more links.")
	}

	target1 += len("a href=") + 1

	target2 := strings.Index(webpage[target1:], "\"")
	target2 = target1 + target2

	link := webpage[target1:target2]

	return link, target2, nil
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
