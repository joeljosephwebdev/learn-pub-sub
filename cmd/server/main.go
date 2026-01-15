package main

import (
	"fmt"
	"log"

	"github.com/joeljosephwebdev/learn-pub-sub/internal/gamelogic"
	"github.com/joeljosephwebdev/learn-pub-sub/internal/pubsub"
	"github.com/joeljosephwebdev/learn-pub-sub/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	// declare a connection string
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	// Call amqp.Dial with the connection string to create a new connection to RabbitMQ.
	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could to connect to RabbitMQ: %v", err)
	}
	// Defer a .Close() of the connection to ensure it's closed when the program exits.
	defer conn.Close()
	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to make channel: %v", err)
	}

	fmt.Println("RabbitMQ connection successfull!")
	gamelogic.PrintServerHelp()

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "pause":
			fmt.Println("Publishing paused game state")
			err = pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: true,
				},
			)
			if err != nil {
				log.Printf("could not publish time: %v", err)
			}
		case "resume":
			fmt.Println("Publishing resumes game state")
			err = pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: false,
				},
			)
			if err != nil {
				log.Printf("could not publish time: %v", err)
			}
		case "quit":
			log.Println("goodbye")
			return
		default:
			fmt.Println("unknown command")
		}
	}

	// Wait for a signal (e.g. Ctrl+C) to exit the program.
	// sigs := make(chan os.Signal, 1)
	// signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// done := make(chan bool, 1)

	// go func() {
	// 	sig := <-sigs
	// 	fmt.Println()
	// 	fmt.Println(sig)
	// 	done <- true
	// }()

	// <-done
	// fmt.Println("Server shutting down...")
}
