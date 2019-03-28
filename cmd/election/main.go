package main

import (
	"fmt"
	"github.com/expectedsh/expected/pkg/election"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	el, err := election.NewElection("test", "localhost:8500", 2*time.Second)
	if err != nil {
		panic(err)
	}

	isLeader := el.ElectLeader(true)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("Job interrupted. Cleaning up")
		err := el.Close()
		if err != nil {
			log.Println("Could not destroy session")
		}
		os.Exit(0)
	}()

	if isLeader {
		fmt.Println("I am leader")
		fmt.Println("Starting to work")
		time.Sleep(17 * time.Second)
		fmt.Println("Work done")
		//panic("LOLILOL")
		err := el.Close()
		if err != nil {
			log.Println("Could not destroy session")
		}
		return
	}

	fmt.Println("I can NOT work. YAY!!!")
}
