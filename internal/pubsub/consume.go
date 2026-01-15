package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType = int

const (
	SimpleQueueDurable SimpleQueueType = iota
	SimpleQueueTransient
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // a simple enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("error creating channel: %v", err)
	}

	queue, err := ch.QueueDeclare(
		queueName,
		queueType == SimpleQueueDurable,
		queueType == SimpleQueueTransient,
		queueType == SimpleQueueTransient,
		false,
		nil,
	)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("error creating queue: %v", err)
	}

	err = ch.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("error binding queue '%s' to exchange '%s': %v", queue.Name, exchange, err)
	}

	return ch, queue, nil
}
