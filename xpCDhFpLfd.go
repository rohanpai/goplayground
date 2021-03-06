package main

import (
	"fmt"
	"github.com/bmizerany/aws4"
	"log"
	"net/http"
	"strings"
)

func main() {
	data := strings.NewReader("{}")
	r, _ := http.NewRequest("POST", "https://dynamodb.us-east-1.amazonaws.com/", data)
	r.Header.Set("Content-Type", "application/x-amz-json-1.0")
	r.Header.Set("X-Amz-Target", "DynamoDB_20111205.ListTables")

	resp, err := aws4.DefaultClient.Do(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.StatusCode)
}
