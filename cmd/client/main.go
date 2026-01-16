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

	gameState := gamelogic.NewGameState(username)

	for {
		words := gamelogic.GetInput()
		if len(words) < 1 {
			continue
		}
		switch words[0] {
		case "spawn":
			err = gameState.CommandSpawn(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "move":
			fmt.Println("we moving")
			armyMove, err := gameState.CommandMove(words)
			if err != nil {
				fmt.Print(err)
				continue
			} else {
				fmt.Printf("move %s %d\n", armyMove.ToLocation, len(armyMove.Units))
			}
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("spamming not allowed yet")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("unknwon command")
		}
	}
}
