package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func main() {
	fmt.Println("Starting Peril server...")
	connectionString := "amqp://guest:guest@localhost:5672/"
	mqConn, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ", err)
	}
	defer mqConn.Close()
	fmt.Println("Connected to Rabbit")

	mqChan, err := mqConn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel", err)
	}
	defer func(mqChan *amqp.Channel) {
		err := mqChan.Close()
		if err != nil {
			log.Fatal("Failed to close channel", err)
		}
	}(mqChan)

	testString := "test this mf"
	err = pubsub.PublishJSON(mqChan, routing.ExchangePerilDirect, routing.PauseKey, testString)
	if err != nil {
		log.Fatal("Failed to publish message", err)
	}

	//signalChan := make(chan os.Signal, 1)
	//signal.Notify(signalChan, os.Interrupt)
	//<-signalChan
}
