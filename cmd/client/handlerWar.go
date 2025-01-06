package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
)

func handlerWar(gs *gamelogic.GameState) func(wo gamelogic.RecognitionOfWar) string {
	return func(wo gamelogic.RecognitionOfWar) string {
		defer fmt.Print("> ")

		outcome, _, _ := gs.HandleWar(wo)
		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return "NackRequeue"
		case gamelogic.WarOutcomeNoUnits:
			return "NackDiscard"
		case gamelogic.WarOutcomeOpponentWon:
			return "Ack"
		case gamelogic.WarOutcomeYouWon:
			return "Ack"
		case gamelogic.WarOutcomeDraw:
			return "Ack"
		default:
			fmt.Println("unknown war outcome")
			return "NackDiscard"
		}
	}
}
