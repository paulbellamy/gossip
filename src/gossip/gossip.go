package gossip

import (
	"errors"
	"fmt"
	. "gossip/message"
	. "gossip/node"
	. "gossip/node_list"
	. "gossip/registry"
	"gossip/util"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

var registry *Registry
var seq uint64

// Start up the RPC Server
func server(registry *Registry, port int) error {
	rpc.Register(registry)
	rpc.HandleHTTP()
	ln, err := net.Listen("tcp", util.Address("", port))
	if err != nil {
		return err
	}
	go func() {
		for {
			conn, _ := ln.Accept()
			rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}()
	return nil
}

// Fetch the initial registry from the address
func connect(registry *Registry, address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	var reply *NodeList
	err = client.Call("Registry.Query", "", &reply)
	if err != nil {
		return err
	}

	client.Close()

	MergeRegistries(registry, reply)
	return nil
}

func client(registry *Registry, seeds []string, port int) error {
	// Get the registries from each of the seeds
	failed := 0
	for _, seed := range seeds {
		if seed != "" {
			if connect(registry, seed) != nil {
				failed++
			}
		}
	}

	// If we couldn't connect to any seeds
	if failed >= len(seeds) {
		return errors.New("Unable to connect to any seed nodes")
	}

	// Announce yourself on the network
	var reply int
	seq++
	message := &Message{
		Origin:        registry.Self,
		Seq:           seq,
		ServiceMethod: "Registry.AddNode",
		Args:          registry.Self,
	}
	return registry.Announce(message, &reply)
}

func Start(name string, hostname string, seeds []string, port int) (chan []byte, error) {
	address := util.Address(hostname, port)
	registry = NewRegistry(name)
	registry.Self = &Node{Name: name, Address: address}

	var err error

	err = server(registry, port)
	if err != nil {
		return nil, err
	}

	err = client(registry, seeds, port)
	if err != nil {
		return nil, err
	}

	return registry.Data, nil
}

func Broadcast(data []byte) {
	var reply int
	seq++
	message := &Message{
		Origin:        registry.Self,
		Seq:           seq,
		ServiceMethod: "Registry.Data",
		Args:          fmt.Sprintf("%s", data),
	}
	registry.Announce(message, &reply)
}
