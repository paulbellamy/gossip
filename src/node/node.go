package main

import (
  "flag"
  "strings"
  "gossip"
  "os"
  "os/signal"
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

  gossip.Start(name, hostname, strings.Split(seeds, ","), port)

  // Wait for Ctrl-C
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)
  <-c
}
