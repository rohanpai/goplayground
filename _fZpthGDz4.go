package main

import (
	&#34;errors&#34;
	&#34;fmt&#34;
	&#34;io/ioutil&#34;
	&#34;net/http&#34;
	&#34;os&#34;
	&#34;strings&#34;
)

const (
	BASEURL string = &#34;http://www.stanford.edu/class/cs193c&#34;
)

func main() {
	resp, err1 := http.Get(BASEURL &#43; &#34;/lectures.html&#34;)
	panicIf(err1)

	defer resp.Body.Close()

	body, err2 := ioutil.ReadAll(resp.Body)
	panicIf(err2)

	bodyString := string(body)

	linkList := make([]string, 0)
	getLinks(bodyString, &amp;linkList)

	fmt.Println(&#34;Fetching Urls...&#34;)

	for _, link := range linkList {
		fmt.Println(link)
	}

	fmt.Println(&#34;Done.&#34;)

	os.Exit(0)
}

func getLinks(webpage string, linkList *[]string) {
	gotLink, i, err := getLink(webpage)
	if err != nil {
		return
	}

	if strings.HasSuffix(gotLink, &#34;.zip&#34;) {
		*linkList = append(*linkList, gotLink)
	}

	if strings.HasSuffix(gotLink, &#34;.pdf&#34;) {
		*linkList = append(*linkList, gotLink)
	}

	getLinks(webpage[i:], linkList)
}

func getLink(webpage string) (string, int, error) {
	target1 := strings.Index(webpage, &#34;a href=\&#34;&#34;)
	if target1 == -1 {
		return &#34;&#34;, -1, errors.New(&#34;No more links.&#34;)
	}

	target1 &#43;= len(&#34;a href=&#34;) &#43; 1

	target2 := strings.Index(webpage[target1:], &#34;\&#34;&#34;)
	target2 = target1 &#43; target2

	link := webpage[target1:target2]

	return link, target2, nil
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
