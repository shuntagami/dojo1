package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/shuntagami/dojo1/kadai3/typing_game"
)

func main() {
	client, err := typing_game.Initialize(os.Getenv("PROJECT_ROOT_DIR"))
	if err != nil {
		log.Fatal(err)
	}

	timeLimit, err := strconv.Atoi(os.Getenv("TIME_LIMIT"))
	if err != nil {
		log.Fatal(err)
	}
	timer1 := time.NewTimer(time.Duration(timeLimit) * time.Second)

	go func() {
		<-timer1.C
		fmt.Println("\nTime out!")
		fmt.Printf("Your score is %d\n", client.Score)
		os.Exit(0)
	}()

	for {
		if err := client.Play(); err != nil {
			log.Fatal(err)
		}
	}
}
