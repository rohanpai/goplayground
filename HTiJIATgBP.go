package main

import (
	&#34;github.com/jjeffery/stomp&#34;
)

func main() {
	// send with receipt and an optional header
	err := c.Send(
		&#34;/queue/test-1&#34;,            // destination
		&#34;text/plain&#34;,               // content-type
		[]byte(&#34;Message number 1&#34;), // body
		stomp.NewHeader(&#34;expires&#34;, &#34;2020-12-31 23:59:59&#34;))
	if err != nil {
		return err
	}

	// send with no receipt and no optional headers
	err = c.Send(&#34;/queue/test-2&#34;, &#34;application/xml&#34;,
		[]byte(&#34;&lt;message&gt;hello&lt;/message&gt;&#34;), nil)
	if err != nil {
		return err
	}

	return nil
}
