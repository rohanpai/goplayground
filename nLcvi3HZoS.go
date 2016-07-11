// This example declares a durable Exchange, an ephemeral (auto-delete) Queue,
// binds the Queue to the Exchange with a binding key, and consumes every
// message published to that Exchange with that routing key.
//
package main

import (
	&#34;code.google.com/p/go.net/websocket&#34;
	&#34;flag&#34;
	&#34;fmt&#34;
	&#34;github.com/streadway/amqp&#34;
	&#34;log&#34;
	&#34;net/http&#34;
)

var (
	uri		= flag.String(&#34;uri&#34;, &#34;amqp://guest:guest@localhost:5672/&#34;, &#34;AMQP URI&#34;)
	exchange	= flag.String(&#34;exchange&#34;, &#34;logs&#34;, &#34;Durable, non-auto-deleted AMQP exchange name&#34;)
	exchangeType	= flag.String(&#34;exchange-type&#34;, &#34;direct&#34;, &#34;Exchange type - direct|fanout|topic|x-custom&#34;)
	queue		= flag.String(&#34;queue&#34;, &#34;logs&#34;, &#34;Ephemeral AMQP queue name&#34;)
	bindingKey	= flag.String(&#34;key&#34;, &#34;&#34;, &#34;AMQP binding key&#34;)
	consumerTag	= flag.String(&#34;consumer-tag&#34;, &#34;simple-consumer&#34;, &#34;AMQP consumer tag (should not be blank)&#34;)
	wsPort		= flag.Int(&#34;ws-port&#34;, 23456, &#34;WebSockets port to listen to&#34;)
)

func init() {
	flag.Parse()
}

func transmitLog(ws *websocket.Conn) {
	// Here, send to the browser what the AMQP Consumer got.
}

func main() {
	err = http.ListenAndServe(fmt.Sprintf(&#34;:%d&#34;, *wsPort), nil)
	if err != nil {
		log.Fatalf(&#34;listen: %s&#34;, err)
	}

	http.Handle(&#34;/&#34;, websocket.Handler(transmitLog))

	c, err := NewConsumer(*uri, *exchange, *exchangeType, *queue, *bindingKey, *consumerTag)
	if err != nil {
		log.Fatalf(&#34;%s&#34;, err)
	}

	select {}

	if err := c.Shutdown(); err != nil {
		log.Fatalf(&#34;error during shutdown: %s&#34;, err)
	}

}

type Consumer struct {
	conn	*amqp.Connection
	channel	*amqp.Channel
	tag	string
	done	chan error
}

func NewConsumer(amqpURI, exchange, exchangeType, queue, key, ctag string) (*Consumer, error) {
	c := &amp;Consumer{
		conn:		nil,
		channel:	nil,
		tag:		ctag,
		done:		make(chan error),
	}

	var err error

	log.Printf(&#34;dialing %s&#34;, amqpURI)
	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf(&#34;Dial: %s&#34;, err)
	}

	log.Printf(&#34;got Connection, getting Channel&#34;)
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf(&#34;Channel: %s&#34;, err)
	}

	log.Printf(&#34;got Channel, declaring Exchange (%s)&#34;, exchange)
	if err = c.channel.ExchangeDeclare(
		exchange,	// name of the exchange
		exchangeType,	// type
		true,		// durable
		false,		// delete when complete
		false,		// internal
		false,		// noWait
		nil,		// arguments
	); err != nil {
		return nil, fmt.Errorf(&#34;Exchange Declare: %s&#34;, err)
	}

	log.Printf(&#34;declared Exchange, declaring Queue (%s)&#34;, queue)
	state, err := c.channel.QueueDeclare(
		queue,	// name of the queue
		true,	// durable
		false,	// delete when usused
		false,	// exclusive
		false,	// noWait
		nil,	// arguments
	)
	if err != nil {
		return nil, fmt.Errorf(&#34;Queue Declare: %s&#34;, err)
	}

	log.Printf(&#34;declared Queue (%d messages, %d consumers), binding to Exchange (key &#39;%s&#39;)&#34;,
		state.Messages, state.Consumers, key)

	if err = c.channel.QueueBind(
		queue,		// name of the queue
		key,		// bindingKey
		exchange,	// sourceExchange
		false,		// noWait
		nil,		// arguments
	); err != nil {
		return nil, fmt.Errorf(&#34;Queue Bind: %s&#34;, err)
	}

	log.Printf(&#34;Queue bound to Exchange, starting Consume (consumer tag &#39;%s&#39;)&#34;, c.tag)
	deliveries, err := c.channel.Consume(
		queue,	// name
		c.tag,	// consumerTag,
		false,	// noAck
		false,	// exclusive
		false,	// noLocal
		false,	// noWait
		nil,	// arguments
	)
	if err != nil {
		return nil, fmt.Errorf(&#34;Queue Consume: %s&#34;, err)
	}

	go handle(deliveries, c.done)

	return c, nil
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf(&#34;Consumer cancel failed: %s&#34;, err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf(&#34;AMQP connection close error: %s&#34;, err)
	}

	defer log.Printf(&#34;AMQP shutdown OK&#34;)

	// wait for handle() to exit
	return &lt;-c.done
}

func handle(deliveries &lt;-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		log.Printf(
			&#34;got %dB delivery: [%v] %s&#34;,
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
	}
	log.Printf(&#34;handle: deliveries channel closed&#34;)
	done &lt;- nil
}
