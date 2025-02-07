package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func main() {
	fmt.Println("Starting Peril server...")
	gamelogic.PrintServerHelp()
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

	_, _, err = pubsub.DeclareAndBind(mqConn, routing.ExchangePerilTopic, routing.GameLogSlug, "game_logs.*", 1)
	if err != nil {
		fmt.Println("Failed to declare and bind game log exchange")
		fmt.Println(err)
		return
	}

	err = pubsub.SubscribeGob(mqConn, routing.ExchangePerilTopic, routing.GameLogSlug, "*", 1, handlerLogs())
	if err != nil {
		log.Fatal("failed to subscribe to logs queue", err)
	}

	serverMessage := routing.PlayingState{
		IsPaused: true,
	}

	for {
		command := gamelogic.GetInput()
		switch command[0] {
		case "pause":

			fmt.Println("pausing")
			serverMessage.IsPaused = true
			err = pubsub.PublishJSON(mqChan, routing.ExchangePerilDirect, routing.PauseKey, serverMessage)
			if err != nil {
				log.Fatal("Failed to publish message", err)
			}
		case "resume":
			fmt.Println("resuming")
			serverMessage.IsPaused = false
			err = pubsub.PublishJSON(mqChan, routing.ExchangePerilDirect, routing.PauseKey, serverMessage)
			if err != nil {
				log.Fatal("Failed to publish message", err)
			}
		case "quit":
			goto loopEnd
		default:
			fmt.Println("unknown command")
		}
	}
loopEnd:
	fmt.Println("Goodbye!")

}
