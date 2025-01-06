package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func handlerWar(gs *gamelogic.GameState) func(wo gamelogic.RecognitionOfWar) string {
	return func(wo gamelogic.RecognitionOfWar) string {
		defer fmt.Print("> ")

		attacker := wo.Attacker.Username
		defender := wo.Defender.Username

		outcome, _, _ := gs.HandleWar(wo)
		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return "NackRequeue"
		case gamelogic.WarOutcomeNoUnits:
			return "NackDiscard"
		case gamelogic.WarOutcomeOpponentWon:
			message := fmt.Sprintf("%s won a war against %s", defender, attacker)
			err := pubsub.PublishGob(gs.MqChan, routing.ExchangePerilTopic, routing.GameLogSlug+"."+attacker, message)
			if err != nil {
				return "NackRequeue"
			}
			return "Ack"
		case gamelogic.WarOutcomeYouWon:
			message := fmt.Sprintf("%s won a war against %s", attacker, defender)
			err := pubsub.PublishGob(gs.MqChan, routing.ExchangePerilTopic, routing.GameLogSlug+"."+attacker, message)
			if err != nil {
				return "NackRequeue"
			}
			return "Ack"
		case gamelogic.WarOutcomeDraw:
			message := fmt.Sprintf("A war between %s and %s resulted in a draw", attacker, defender)
			err := pubsub.PublishGob(gs.MqChan, routing.ExchangePerilTopic, routing.GameLogSlug+"."+attacker, message)
			if err != nil {
				return "NackRequeue"
			}
			return "Ack"
		default:
			fmt.Println("unknown war outcome")
			return "NackDiscard"
		}
	}
}
