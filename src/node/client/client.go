package client

import (
  "fmt"
  "net/rpc"
  . "node/registry"
  "node/util"
)

func Start(name string, seeds []string, port int) {
  for _, seed := range seeds {
    address := util.Address(seed, port)
    client, err := rpc.DialHTTP("tcp", address)
    if err != nil {
      fmt.Println("Error Connecting to Seed", seed, ":", err)
    } else {
      fmt.Println("Connected to Seed :", seed)
      // Asynchronous call
      var reply *int
      node := &Node{Name: name}
      divCall := client.Go("Registry.AddNode", node, &reply, nil)
      <-divCall.Done// will be equal to divCall
      // check errors, print, etc.
    }
  }
}
