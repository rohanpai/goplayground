package main

import (
	&#34;fmt&#34;
	&#34;github.com/bmizerany/aws4&#34;
	&#34;log&#34;
	&#34;net/http&#34;
	&#34;strings&#34;
)

func main() {
	data := strings.NewReader(&#34;{}&#34;)
	r, _ := http.NewRequest(&#34;POST&#34;, &#34;https://dynamodb.us-east-1.amazonaws.com/&#34;, data)
	r.Header.Set(&#34;Content-Type&#34;, &#34;application/x-amz-json-1.0&#34;)
	r.Header.Set(&#34;X-Amz-Target&#34;, &#34;DynamoDB_20111205.ListTables&#34;)

	resp, err := aws4.DefaultClient.Do(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.StatusCode)
}
