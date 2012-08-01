package main

import (
	"flag"
	"gossip"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
)

// Flags
var name string
var seeds string
var port int

func main() {
	hostname, _ := os.Hostname()
	flag.StringVar(&name, "name", hostname, "Name of this node, must be unique across the cluster")
	flag.StringVar(&seeds, "seeds", "", "Comma-seperated list of seed-nodes")
	flag.IntVar(&port, "port", 5000, "Port to listen on")
	flag.Parse()

	received, err := gossip.Start(name, hostname, strings.Split(seeds, ","), port)
	if err != nil {
		log.Println(err)
		return
	} else {
		log.Println("Client started")
	}

	// Print any messages received
	go func() {
		for {
			d := <-received
			log.Printf("Received: %s", d)
		}
	}()

	// Send a dummy message out every 5 seconds to test we can broadcast
	go func() {
		for {
			time.Sleep(5 * time.Second)
			log.Printf("Announced: %s", name)
			gossip.Broadcast([]byte(name))
		}
	}()

	// Wait for Ctrl-C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
