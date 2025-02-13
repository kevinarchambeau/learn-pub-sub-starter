package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"strconv"
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
	gs.MqChan = mqChan

	_, _, err = pubsub.DeclareAndBind(mqConn, routing.ExchangePerilDirect, "pause."+userName, routing.PauseKey, 0)
	if err != nil {
		return
	}

	_, _, err = pubsub.DeclareAndBind(mqConn, routing.ExchangePerilTopic, routing.ArmyMovesPrefix+"."+userName, routing.ArmyMovesPrefix+".*", 0)
	if err != nil {
		return
	}

	// handle server pauses
	err = pubsub.SubscribeJSON(mqConn, routing.ExchangePerilDirect, "pause."+userName, routing.PauseKey, 0, handlerPause(gs))
	if err != nil {
		log.Fatal("failed to subscribe to pause queue", err)
	}
	//handle other player moves
	err = pubsub.SubscribeJSON(mqConn, routing.ExchangePerilTopic, routing.ArmyMovesPrefix+"."+userName, routing.ArmyMovesPrefix+".*", 0, handlerMove(gs))
	if err != nil {
		log.Fatal("failed to subscribe to move queue", err)
	}

	// handle wars
	err = pubsub.SubscribeJSON(mqConn, routing.ExchangePerilTopic, routing.WarRecognitionsPrefix, routing.WarRecognitionsPrefix+".*", 1, handlerWar(gs))
	if err != nil {
		log.Fatal("failed to subscribe to war queue", err)
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
			move, err := gs.CommandMove(command)
			if err != nil {
				fmt.Println(err)
			}
			err = pubsub.PublishJSON(mqChan, routing.ExchangePerilTopic, routing.ArmyMovesPrefix+"."+userName, move)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("move published")

		case "spam":
			arg, err := strconv.Atoi(command[1])
			if err != nil {
				fmt.Println("invalid argument")
				break
			}
			for i := 0; i < arg; i++ {
				message := gamelogic.GetMaliciousLog()

				err := pubsub.PublishGob(gs.MqChan, routing.ExchangePerilTopic, routing.GameLogSlug+"."+userName, message)
				if err != nil {
					fmt.Println("error publishing spam log")
				}
			}
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
