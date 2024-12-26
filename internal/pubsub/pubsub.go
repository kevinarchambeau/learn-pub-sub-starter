package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	jsonString, err := json.Marshal(val)
	if err != nil {
		log.Printf("failed to marshal json: %s", err)
	}

	err = ch.PublishWithContext(context.Background(), exchange, key, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonString,
		},
	)
	if err != nil {
		log.Fatal("publish failed", err)
	}
	return nil
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	mqChan, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("failed to open a channel: %s", err)
	}

	isDurable := simpleQueueType == 1
	isTransient := simpleQueueType == 0
	queue, err := mqChan.QueueDeclare(queueName, isDurable, isTransient, isTransient, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("failed to create queue: %s", err)
	}
	err = mqChan.QueueBind(queue.Name, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("failed to bind to exchange: %s", err)
	}

	return mqChan, queue, nil
}
