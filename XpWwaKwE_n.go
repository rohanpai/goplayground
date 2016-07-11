package main

import (
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"testing"
)

import (
	nsq "github.com/bitly/go-nsq"
	"github.com/bitly/nsq/nsqd"
)

type nopLogger struct{}

func (*nopLogger) Output(int, string) error {
	return nil
}

func newDaemon() *nsqd.NSQD {
	opts := nsqd.NewNSQDOptions()

	// Disable http/https
	opts.HTTPAddress = ""
	opts.HTTPSAddress = ""
	// Disable logging
	opts.Logger = &nopLogger{}
	// Do not create on disc queue
	opts.DataPath = "/dev/null"

	nsqd := nsqd.NewNSQD(opts)
	nsqd.Main()
	return nsqd
}

type consumer struct{ *nsq.Consumer }

func (c *consumer) Stop() {
	c.Consumer.Stop()
	<-c.Consumer.StopChan
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
	r.SetLogger(&nopLogger{}, 0)

	// Set the handler
	r.AddHandler(hdlr)

	// Connect to the NSQ daemon
	if err := r.ConnectToNSQD(tcpAddr); err != nil {
		t.Fatal(err)
	}

	return &consumer{Consumer: r}
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
	p.SetLogger(&nopLogger{}, 0)

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
	hdlr := func(msg *nsq.Message) error { msgs <- msg.Body; return nil }
	consumer := newConsumer(b, "localhost:4150", "mytopic", "mychan1", hdlr)
	defer consumer.Stop()

	// Create producer
	producer := newProducer(b, "localhost:4150")
	defer producer.Stop()

	// Tell Go to use as many cores as available
	b.SetParallelism(runtime.NumCPU())

	// reset Go's timer.
	b.ResetTimer()

	// Run in Parallel
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Send "hello world" to NSQ and wait for it to arrive.
			if err := producer.Publish("mytopic", []byte("hello world")); err != nil {
				b.Fatal(err)
			}
			if msg, ok := <-msgs; !ok {
				b.Fatal("Message chan closed.")
			} else if expect, got := "hello world", string(msg); expect != got {
				b.Fatalf("Unexpected message. Expected: %s, Got: %s", expect, got)
			}
		}
	})
}
