package main

import (
	&#34;io/ioutil&#34;
	&#34;log&#34;
	&#34;os&#34;
	&#34;runtime&#34;
	&#34;testing&#34;
)

import (
	nsq &#34;github.com/bitly/go-nsq&#34;
	&#34;github.com/bitly/nsq/nsqd&#34;
)

type nopLogger struct{}

func (*nopLogger) Output(int, string) error {
	return nil
}

func newDaemon() *nsqd.NSQD {
	opts := nsqd.NewNSQDOptions()

	// Disable http/https
	opts.HTTPAddress = &#34;&#34;
	opts.HTTPSAddress = &#34;&#34;
	// Disable logging
	opts.Logger = &amp;nopLogger{}
	// Do not create on disc queue
	opts.DataPath = &#34;/dev/null&#34;

	nsqd := nsqd.NewNSQD(opts)
	nsqd.Main()
	return nsqd
}

type consumer struct{ *nsq.Consumer }

func (c *consumer) Stop() {
	c.Consumer.Stop()
	&lt;-c.Consumer.StopChan
}

func newConsumer(t testing.TB, tcpAddr, topicName, channelName string, hdlr nsq.HandlerFunc) *consumer {
	// Create the configuration object and set the maxInFlight
	cfg := nsq.NewConfig()
	cfg.MaxInFlight = 8

	// Create the consumer with the given topic and chanel names
	r, err := nsq.NewConsumer(topicName, channelName, cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Disable logging
	r.SetLogger(&amp;nopLogger{}, 0)

	// Set the handler
	r.AddHandler(hdlr)

	// Connect to the NSQ daemon
	if err := r.ConnectToNSQD(tcpAddr); err != nil {
		t.Fatal(err)
	}

	return &amp;consumer{Consumer: r}
}

func newProducer(t testing.TB, tcpAddr string) *nsq.Producer {
	// Create the configuration object and set the maxInFlight
	cfg := nsq.NewConfig()
	cfg.MaxInFlight = 8

	// Create the producer
	p, err := nsq.NewProducer(tcpAddr, cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Disable logging
	p.SetLogger(&amp;nopLogger{}, 0)

	return p
}

func BenchmarkPubSub(b *testing.B) {
	// Disable general logging
	log.SetOutput(ioutil.Discard)
	defer func() { log.SetOutput(os.Stderr) }()

	// Start NSQD and make sure to shut it down when leaving
	nsqd := newDaemon()
	defer nsqd.Exit()

	// Create the consumer and send every message to the chan
	msgs := make(chan []byte)
	hdlr := func(msg *nsq.Message) error { msgs &lt;- msg.Body; return nil }
	consumer := newConsumer(b, &#34;localhost:4150&#34;, &#34;mytopic&#34;, &#34;mychan1&#34;, hdlr)
	defer consumer.Stop()

	// Create producer
	producer := newProducer(b, &#34;localhost:4150&#34;)
	defer producer.Stop()

	// Tell Go to use as many cores as available
	b.SetParallelism(runtime.NumCPU())

	// reset Go&#39;s timer.
	b.ResetTimer()

	// Run in Parallel
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Send &#34;hello world&#34; to NSQ and wait for it to arrive.
			if err := producer.Publish(&#34;mytopic&#34;, []byte(&#34;hello world&#34;)); err != nil {
				b.Fatal(err)
			}
			if msg, ok := &lt;-msgs; !ok {
				b.Fatal(&#34;Message chan closed.&#34;)
			} else if expect, got := &#34;hello world&#34;, string(msg); expect != got {
				b.Fatalf(&#34;Unexpected message. Expected: %s, Got: %s&#34;, expect, got)
			}
		}
	})
}
