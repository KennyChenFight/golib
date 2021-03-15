package amqplib

import (
	"github.com/assembla/cony"
	"github.com/streadway/amqp"
)

type AMQPPublisher struct {
	Publisher *cony.Publisher
}

func (p *AMQPPublisher) Publish(data amqp.Publishing) (err error) {
	return p.Publisher.Publish(data)
}

func (p *AMQPPublisher) Close() {
	p.Publisher.Cancel()
}
