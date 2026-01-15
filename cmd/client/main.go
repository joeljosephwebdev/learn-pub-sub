package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joeljosephwebdev/learn-pub-sub/internal/gamelogic"
	"github.com/joeljosephwebdev/learn-pub-sub/internal/pubsub"
	"github.com/joeljosephwebdev/learn-pub-sub/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	// declare a connection string
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("error connecting to rabbitmq: %v", err)
	}
	defer conn.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatal(err)
	}

	_, queue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		fmt.Sprintf("%s.%s", routing.PauseKey, username),
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
	)
	if err != nil {
		log.Fatalf("failed to declare and bind queue: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	// Wait for a signal (e.g. Ctrl+C) to exit the program.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	<-done
	fmt.Println("RabbitMQ connection closed.")
}
