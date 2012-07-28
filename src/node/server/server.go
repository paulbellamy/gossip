package server

import (
  "fmt"
  "net"
  "net/http"
  "net/rpc"
  . "node/node"
  . "node/registry"
  "node/util"
)

func Start(name string, hostname string, port int) {
  address := util.Address(hostname, port)
  registry := NewRegistry(name)
  registry.Self = &Node{Name: name, Address: address}
  rpc.Register(registry)
  rpc.HandleHTTP()
  ln, err := net.Listen("tcp", util.Address("", port))
  if err != nil {
    fmt.Println("Failed to start server on port", port,":",err)
  } else {
    fmt.Println("Server Listening on port", port)
  }
  http.Serve(ln, nil)
}
