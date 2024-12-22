package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"os/signal"
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
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
