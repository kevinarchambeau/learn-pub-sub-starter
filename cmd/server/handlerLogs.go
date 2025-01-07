package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"time"
)

func handlerLogs() func(message string) string {
	return func(message string) string {
		defer fmt.Print("> ")
		log := routing.GameLog{
			CurrentTime: time.Time{},
			Message:     message,
			Username:    "server",
		}
		err := gamelogic.WriteLog(log)
		if err != nil {
			return "NackRequeue"
		}
		return "Ack"
	}
}
