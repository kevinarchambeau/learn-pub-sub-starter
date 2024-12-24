package pubsub

import (
	"context"
	"encoding/json"
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
