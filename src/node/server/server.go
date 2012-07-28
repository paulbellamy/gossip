package server

import (
  "fmt"
  "net"
  "net/http"
  "node/util"
  "net/rpc"
  . "node/registry"
)

func Start(name, listen string, port int) {
  address := util.Address(listen, port)

  registry := NewRegistry(name)
  rpc.Register(registry)
  rpc.HandleHTTP()
  ln, err := net.Listen("tcp", address)
  if err != nil {
    fmt.Println("Failed to start server on", address,":",err)
  } else {
    fmt.Println("Server Listening on", address)
  }
  http.Serve(ln, nil)
}
