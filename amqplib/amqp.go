package amqplib

import (
	"errors"
	"github.com/assembla/cony"
	"github.com/streadway/amqp"
	"log"
)

type AMQPConnectionConfig struct {
	URL          string
	ErrorHandler func(err error)
}

type ExchangeType string

const (
	Fanout ExchangeType = "fanout"
	Direct ExchangeType = "direct"
	Topic  ExchangeType = "topic"
)

type AMQPQueueConfig struct {
	ExchangeName string
	ExchangeType ExchangeType
	RoutingKey   string
	QueueName    string
	ConsumerName string
	AutoDelete   bool
}

func NewAMQPClient(connConf *AMQPConnectionConfig) *AMQPClient {
	client := cony.NewClient(cony.URL(connConf.URL))
	errHandler := connConf.ErrorHandler

	if errHandler == nil {
		errHandler = func(err error) {
			log.Println("amqp error:", err)
		}
	}

	go func() {
		for client.Loop() {
			select {
			case err := <-client.Errors():
				errHandler(err)
			case b := <-client.Blocking():
				if b.Active {
					log.Println("the amqp connection is blocked by amqp server:", b.Reason)
				} else {
					log.Println("unknown error for blocking")
				}
			}
		}
	}()

	c := &AMQPClient{
		client:           client,
		connectionConfig: connConf,
		errorHandler:     errHandler,
	}
	return c
}

type AMQPClient struct {
	client           *cony.Client
	connectionConfig *AMQPConnectionConfig
	errorHandler     func(err error)
}

func (c *AMQPClient) NewPublisher(cfg *AMQPQueueConfig) (*AMQPPublisher, error) {
	exchange, err := c.initExchange(cfg)
	if err != nil {
		return nil, err
	}

	routingKey, err := c.initRoutingKey(cfg)
	if err != nil {
		return nil, err
	}

	queue, err := c.initQueue(cfg)
	if err != nil {
		return nil, err
	}

	publisher := cony.NewPublisher(exchange.Name, routingKey)

	c.client.Declare(c.initDeclare(exchange, queue, routingKey))
	c.client.Publish(publisher)
	return &AMQPPublisher{
		Publisher: publisher,
	}, nil
}

func (c *AMQPClient) NewConsumer(cfg *AMQPQueueConfig) (*AMQPConsumer, error) {
	exchange, err := c.initExchange(cfg)
	if err != nil {
		return nil, err
	}

	routingKey, err := c.initRoutingKey(cfg)
	if err != nil {
		return nil, err
	}

	queue, err := c.initQueue(cfg)
	if err != nil {
		return nil, err
	}

	consumer := cony.NewConsumer(queue)
	c.client.Declare(c.initDeclare(exchange, queue, routingKey))
	c.client.Consume(consumer)
	return &AMQPConsumer{
		amqpClient:  c,
		consumer:    consumer,
		deliveryCh:  make(chan amqp.Delivery),
		queueConfig: cfg,
		closed:      false,
	}, nil
}

func (c *AMQPClient) Close() {
	c.client.Close()
}

func (c *AMQPClient) initExchange(cfg *AMQPQueueConfig) (*cony.Exchange, error) {
	if cfg.ExchangeName == "" {
		return nil, errors.New("exchange name is required")
	}
	if cfg.ExchangeType == "" {
		return nil, errors.New("exchange type is required")
	}
	return &cony.Exchange{
		Name:       cfg.ExchangeName,
		Kind:       string(cfg.ExchangeType),
		Durable:    true,
		AutoDelete: cfg.AutoDelete,
	}, nil
}

func (c *AMQPClient) initRoutingKey(cfg *AMQPQueueConfig) (string, error) {
	if cfg.ExchangeType != Fanout && cfg.RoutingKey == "" {
		return "", errors.New("routing key is required")
	}
	return cfg.RoutingKey, nil
}

func (c *AMQPClient) initQueue(cfg *AMQPQueueConfig) (*cony.Queue, error) {
	name := cfg.QueueName
	if name == "" {
		return nil, errors.New("queue name is required")
	}

	return &cony.Queue{
		Name:       name,
		Durable:    true,
		AutoDelete: cfg.AutoDelete,
	}, nil
}

func (c *AMQPClient) initDeclare(exchange *cony.Exchange, queue *cony.Queue, routingKey string) []cony.Declaration {
	declarations := []cony.Declaration{
		cony.DeclareExchange(*exchange),
		cony.DeclareQueue(queue),
		cony.DeclareBinding(cony.Binding{
			Queue:    queue,
			Exchange: *exchange,
			Key:      routingKey,
		}),
	}
	return declarations
}
