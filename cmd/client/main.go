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
	fmt.Println("Starting Peril client...")
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

	userName, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatal("Failed to get username", err)
	}
	gs := gamelogic.NewGameState(userName)

	_, _, err = pubsub.DeclareAndBind(mqConn, routing.ExchangePerilDirect, "pause."+userName, routing.PauseKey, 0)
	if err != nil {
		return
	}

	for {
		command := gamelogic.GetInput()
		switch command[0] {
		case "status":
			gs.CommandStatus()
		case "spawn":
			err = gs.CommandSpawn(command)
			if err != nil {
				fmt.Println(err)
			}
		case "move":
			_, err = gs.CommandMove(command)
			if err != nil {
				fmt.Println(err)
			}
		case "spam":
			fmt.Println("Not supported yet")
		case "help":
			gamelogic.PrintClientHelp()
		case "quit":
			goto loopEnd
		default:
			fmt.Println("unknown command")
		}
	}
loopEnd:
	gamelogic.PrintQuit()

}
