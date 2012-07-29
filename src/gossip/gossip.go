package gossip

import (
  "encoding/gob"
  "log"
  . "gossip/message"
  . "gossip/node"
  . "gossip/registry"
  "gossip/util"
  "net"
  "net/http"
  "net/rpc"
)

// Start up the RPC Server
func server(registry *Registry, port int) {
  rpc.Register(registry)
  rpc.HandleHTTP()
  ln, err := net.Listen("tcp", util.Address("", port))
  if err != nil {
    log.Println("Failed to start server on port", port,":",err)
  } else {
    log.Println("Server Listening on port", port)
  }
  http.Serve(ln, nil)
}

// Fetch the initial registry from the address
func connect(registry *Registry, address string) error {
  client, err := rpc.DialHTTP("tcp", address)
  if err != nil {
    return err
  }

  var reply *Registry
  err = client.Call("Registry.Query", "", &reply)
  if err != nil {
    return err
  }

  MergeRegistries(registry, reply)
  return nil
}

func client(registry *Registry, seeds []string, port int) {
  var err error
  var seq uint64 // sent message sequence count

  // Get the registries from each of the seeds
  for _, seed := range seeds {
    err = connect(registry, seed)
    if err != nil {
      log.Println("Error connecting to seed",seed,":",err)
      err = nil
    }
  }

  // Announce yourself on the network
  var reply int
  seq++
  message := &Message{
    Origin: registry.Self,
    Seq: seq,
    ServiceMethod: "Registry.AddNode",
    Args: registry.Self,
  }
  gob.Register(registry.Self)
  registry.Announce(message, &reply)

  log.Println("Registry:",*registry)
}

func Start(name string, hostname string, seeds []string, port int) {
  address := util.Address(hostname, port)
  registry := NewRegistry(name)
  registry.Self = &Node{Name: name, Address: address}

  go server(registry, port)
  go client(registry, seeds, port)
}
