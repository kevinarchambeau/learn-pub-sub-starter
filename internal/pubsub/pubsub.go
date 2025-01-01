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

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int,
	handler func(T) string,
) error {
	mqChan, queue, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return err
	}
	messages, err := mqChan.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		err := func() error {
			for message := range messages {
				var val T
				err = json.Unmarshal(message.Body, &val)
				if err != nil {

					return err
				}
				ackType := handler(val)
				switch ackType {
				case "Ack":
					err = message.Ack(false)
					fmt.Println("Ack")
				case "NackDiscard":
					err = message.Nack(false, false)
					fmt.Println("NackDiscard")
				case "NackRequeue":
					err = message.Nack(false, true)
					fmt.Println("NackRequeue")
				}
				if err != nil {
					return err
				}
			}
			return nil
		}()
		if err != nil {
		}
	}()

	return nil
}
