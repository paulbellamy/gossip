package main

import (
  "flag"
  "strings"
  "node/server"
  "node/client"
  "os"
)


// Flags
var name string
var seeds string
var port int
var listen string
func init() {
    name, _ := os.Hostname()
    flag.StringVar(&name, "name", name, "Name of this node, must be unique across the cluster")
    flag.StringVar(&seeds, "seeds", "127.0.0.1", "Comma-seperated list of seed-nodes")
    flag.IntVar(&port, "port", 5000, "Port which the cluster is on")
    flag.StringVar(&listen, "listen", "", "Binding address for the server to listen on")
}

func main() {
  flag.Parse()

  client.Start(name, strings.Split(seeds, ","), port)
  server.Start(name, listen, port)
}
