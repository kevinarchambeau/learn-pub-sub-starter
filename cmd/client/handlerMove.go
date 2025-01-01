package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
)

func handlerMove(gs *gamelogic.GameState) func(move gamelogic.ArmyMove) string {
	return func(move gamelogic.ArmyMove) string {
		defer fmt.Print("> ")
		outcome := gs.HandleMove(move)
		switch outcome {
		case gamelogic.MoveOutComeSafe:
			return "Ack"
		case gamelogic.MoveOutcomeMakeWar:
			return "Ack"
		case gamelogic.MoveOutcomeSamePlayer:
			return "NackDiscard"
		default:
			return "NackDiscard"
		}
	}
}
