package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func handlerMove(gs *gamelogic.GameState) func(move gamelogic.ArmyMove) string {
	return func(move gamelogic.ArmyMove) string {
		defer fmt.Print("> ")
		outcome := gs.HandleMove(move)
		switch outcome {
		case gamelogic.MoveOutComeSafe:
			return "Ack"
		case gamelogic.MoveOutcomeMakeWar:
			message := gamelogic.RecognitionOfWar{
				Attacker: move.Player,
				Defender: gs.GetPlayerSnap(),
			}
			err := pubsub.PublishJSON(gs.MqChan, routing.ExchangePerilTopic, routing.WarRecognitionsPrefix+"."+gs.GetUsername(), message)
			if err != nil {
				fmt.Println("Failed to publish war recognition")
				return "NackRequeue"
			}
			return "Ack"
		case gamelogic.MoveOutcomeSamePlayer:
			return "Ack"
		default:
			return "NackDiscard"
		}
	}
}
