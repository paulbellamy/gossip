package util

import (
  "fmt"
  "net/rpc"
  . "node/node"
)

func Address(hostname string, port int) string {
  return fmt.Sprintf("%s:%d", hostname, port)
}

func Publish(node *Node, address string) error {
  client, err := rpc.DialHTTP("tcp", address)
  if err != nil {
    return err
  }

  var reply *int
  divCall := client.Go("Registry.AddNode", node, &reply, nil)
  <-divCall.Done// will be equal to divCall
  return nil
}
