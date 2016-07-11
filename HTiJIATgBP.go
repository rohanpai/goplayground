package main

import (
	"github.com/jjeffery/stomp"
)

func main() {
	// send with receipt and an optional header
	err := c.Send(
		"/queue/test-1",            // destination
		"text/plain",               // content-type
		[]byte("Message number 1"), // body
		stomp.NewHeader("expires", "2020-12-31 23:59:59"))
	if err != nil {
		return err
	}

	// send with no receipt and no optional headers
	err = c.Send("/queue/test-2", "application/xml",
		[]byte("<message>hello</message>"), nil)
	if err != nil {
		return err
	}

	return nil
}
