package main

import (
	&#34;fmt&#34;
	&#34;math/rand&#34;
	&#34;time&#34;
)

type Metric string

// BulkUploadMessages will batch up to 10 messages from ch and send
// them to upload().  Rather than block for all 10 messages, it will
// call upload() directly with any number of Metrics if ch is empty.
func BulkUploadMessages(ch &lt;-chan Metric) {
	maximumItemsPerPost := 10
	bulkPost := make([]Metric, 0, maximumItemsPerPost)
	for metric := range ch {
		bulkPost = append(bulkPost[:0], metric)
	outer:
		for len(bulkPost) &lt; maximumItemsPerPost {
			select {
			case metric, ok := &lt;-ch:
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
	fmt.Printf(&#34;HTTP POST len=%d\n&#34;, len(bulkPost))
}

// ProduceMetrics sends &#34;hello world&#34; 50 times to ch with
// sleeps between random Metrics.
func ProduceMetrics(ch chan&lt;- Metric) {
	defer close(ch)
	for i := 0; i &lt; 50; i&#43;&#43; {
		ch &lt;- Metric(&#34;hello world&#34;)
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
