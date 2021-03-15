// Package amqplib is for encapsulating github.com/assembla/cony any operations
//
// As a quick start for publisher:
// 	connectionConfig := &amqplib.AMQPConnectionConfig{
//		URL:          "amqp://guest:guest@localhost/",
//		ErrorHandler: nil,
//	}
//	queueConfig := &amqplib.AMQPQueueConfig{
//		ExchangeName: "test-exchange",
//		ExchangeType: amqplib.Fanout,
//		QueueName:    "test-queue",
//		AutoDelete:   false,
//	}
//	client1 := amqplib.NewAMQPClient(connectionConfig)
//	defer client1.Close()
//	publisher, err := client1.NewPublisher(queueConfig)
//	if err != nil {
//		panic(err)
//	}
//	defer publisher.Close()
//  err := publisher.Publish(amqp.Publishing{
//			Body: []byte("fuck you"),
//		})
//	if err != nil {
//		panic(err)
//	}
//
// As a quick start for consumer:
// 	connectionConfig := &amqplib.AMQPConnectionConfig{
//		URL:          "amqp://guest:guest@localhost/",
//		ErrorHandler: nil,
//	}
//	queueConfig := &amqplib.AMQPQueueConfig{
//		ExchangeName: "test-exchange",
//		ExchangeType: amqplib.Fanout,
//		QueueName:    "test-queue",
//		AutoDelete:   false,
//	}
//	client := amqplib.NewAMQPClient(connectionConfig)
//	defer client1.Close()
//	consumer, err := client.NewConsumer(queueConfig)
//	if err != nil {
//		panic(err)
//	}
//	defer consumer.Close()
// 	for delivery := range consumer.Consume() {
//		fmt.Println(string(delivery.Body))
//		delivery.Ack(false)
//	}
package amqplib

import "github.com/streadway/amqp"

type Publisher interface {
	Publish(data amqp.Publishing) error
	Close()
}

type Consumer interface {
	Consume() <-chan amqp.Delivery
	Close()
}
