package main

import (
	&#34;github.com/streadway/amqp&#34;
	&#34;log&#34;
)

func main() {
	// This example acts as a bridge, shoveling all messages sent from the source
	// exchange &#34;log&#34; to destination exchange &#34;log&#34;.

	// Confirming publishes can help from overproduction and ensure every message
	// is delivered.

	// Setup the source of the store and forward
	source, err := amqp.Dial(&#34;amqp://source/&#34;)
	if err != nil {
		log.Fatalf(&#34;connection.open source: %s&#34;, err)
	}
	defer source.Close()

	chs, err := source.Channel()
	if err != nil {
		log.Fatalf(&#34;channel.open source: %s&#34;, err)
	}

	if err := chs.ExchangeDeclare(&#34;log&#34;, &#34;topic&#34;, true, false, false, false, nil); err != nil {
		log.Fatalf(&#34;exchange.declare destination: %s&#34;, err)
	}

	if _, err := chs.QueueDeclare(&#34;remote-tee&#34;, true, true, false, false, nil); err != nil {
		log.Fatalf(&#34;queue.declare source: %s&#34;, err)
	}

	if err := chs.QueueBind(&#34;remote-tee&#34;, &#34;#&#34;, &#34;logs&#34;, false, nil); err != nil {
		log.Fatalf(&#34;queue.bind source: %s&#34;, err)
	}

	shovel, err := chs.Consume(&#34;remote-tee&#34;, &#34;shovel&#34;, false, false, false, false, nil)
	if err != nil {
		log.Fatalf(&#34;basic.consume source: %s&#34;, err)
	}

	// Setup the destination of the store and forward
	destination, err := amqp.Dial(&#34;amqp://destination/&#34;)
	if err != nil {
		log.Fatalf(&#34;connection.open destination: %s&#34;, err)
	}
	defer destination.Close()

	chd, err := destination.Channel()
	if err != nil {
		log.Fatalf(&#34;channel.open destination: %s&#34;, err)
	}

	if err := chd.ExchangeDeclare(&#34;log&#34;, &#34;topic&#34;, true, false, false, false, nil); err != nil {
		log.Fatalf(&#34;exchange.declare destination: %s&#34;, err)
	}

	// Buffer of 1 for our single outstanding publishing
	pubAcks, pubNacks := chd.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))

	if err := chd.Confirm(false); err != nil {
		log.Fatalf(&#34;confirm.select destination: %s&#34;, err)
	}

	// Now pump the messages, one by one, a smarter implementation
	// would batch the deliveries and use multiple ack/nacks
	for {
		msg, ok := &lt;-shovel
		if !ok {
			log.Fatalf(&#34;source channel closed, see the reconnect example for handling this&#34;)
		}

		err = chd.Publish(&#34;logs&#34;, msg.RoutingKey, false, false, amqp.Publishing{
			// Copy all the properties
			ContentType:     msg.ContentType,
			ContentEncoding: msg.ContentEncoding,
			DeliveryMode:    msg.DeliveryMode,
			Priority:        msg.Priority,
			CorrelationId:   msg.CorrelationId,
			ReplyTo:         msg.ReplyTo,
			Expiration:      msg.Expiration,
			MessageId:       msg.MessageId,
			Timestamp:       msg.Timestamp,
			Type:            msg.Type,
			UserId:          msg.UserId,
			AppId:           msg.AppId,

			// Custom headers
			Headers: msg.Headers,

			// And the body
			Body: msg.Body,
		})

		if err != nil {
			msg.Nack(false, false)
			log.Fatalf(&#34;basic.publish destination: %s&#34;, msg)
		}

		// only ack the source delivery when the destination acks the publishing
		// here you could check for delivery order by keeping a local state of
		// expected delivery tags
		select {
		case &lt;-pubAcks:
			msg.Ack(false)
		case &lt;-pubNacks:
			msg.Nack(false, false)
		}
	}
}
