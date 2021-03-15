package amqplib

import (
	"github.com/assembla/cony"
	"github.com/streadway/amqp"
)

type AMQPConsumer struct {
	amqpClient  *AMQPClient
	consumer    *cony.Consumer
	deliveryCh  chan amqp.Delivery
	queueConfig *AMQPQueueConfig
	closed      bool
}

func (c *AMQPConsumer) Consume() <-chan amqp.Delivery {
	go func() {
		for !c.closed {
			select {
			case msg := <-c.consumer.Deliveries():
				c.deliveryCh <- msg
			case err := <-c.consumer.Errors():
				c.amqpClient.errorHandler(err)
			}
		}
	}()
	return c.deliveryCh
}

func (c *AMQPConsumer) Close() {
	if !c.closed {
		c.closed = true
		c.consumer.Cancel()
	}

	if c.deliveryCh != nil {
		close(c.deliveryCh)
		c.deliveryCh = nil
	}
}
