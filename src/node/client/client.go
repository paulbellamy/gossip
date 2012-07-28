package client

import (
  "fmt"
  "net/rpc"
  . "node/node"
  "node/util"
)

func Connect(node *Node, address string) error {
  client, err := rpc.DialHTTP("tcp", address)
  if err != nil {
    return err
  }

  var reply *int
  divCall := client.Go("Registry.QueryAll", node.Address, &reply, nil)
  <-divCall.Done// will be equal to divCall

  util.Publish(node, address)
  return nil
}

func Start(name string, hostname string, seeds []string, port int) {
  var err error

  address := util.Address(hostname, port)

  node := &Node{Name: name, Address: address}

  for _, seed := range seeds {
    err = Connect(node, seed)
    if err != nil {
      fmt.Println("Error connecting to seed",seed,":",err)
      err = nil
    }
  }
}
