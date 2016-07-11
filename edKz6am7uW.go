package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Metric string

// BulkUploadMessages will batch up to 10 messages from ch and send
// them to upload().  Rather than block for all 10 messages, it will
// call upload() directly with any number of Metrics if ch is empty.
func BulkUploadMessages(ch <-chan Metric) {
	maximumItemsPerPost := 10
	bulkPost := make([]Metric, 0, maximumItemsPerPost)
	for metric := range ch {
		bulkPost = append(bulkPost[:0], metric)
	outer:
		for len(bulkPost) < maximumItemsPerPost {
			select {
			case metric, ok := <-ch:
				if !ok {
					break outer
				}
				bulkPost = append(bulkPost, metric)
			default:
				break outer
			}
		}
		upload(bulkPost)
	}
}

func upload(bulkPost []Metric) {
	fmt.Printf("HTTP POST len=%d\n", len(bulkPost))
}

// ProduceMetrics sends "hello world" 50 times to ch with
// sleeps between random Metrics.
func ProduceMetrics(ch chan<- Metric) {
	defer close(ch)
	for i := 0; i < 50; i++ {
		ch <- Metric("hello world")
		if rand.Intn(10) == 0 {
			time.Sleep(time.Nanosecond)
		}
	}
}

func main() {
	ch := make(chan Metric, 1024)
	go ProduceMetrics(ch)
	BulkUploadMessages(ch)
}
