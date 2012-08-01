package gossip

import (
	"encoding/gob"
	"errors"
	. "gossip/message"
	. "gossip/node"
	. "gossip/node_list"
	. "gossip/registry"
	"gossip/util"
	"net/http"
	"net/rpc"
)

var registry *Registry
var seq uint64

// Start up the RPC Server
func server(registry *Registry, port int) error {
	rpc.Register(registry)
	rpc.HandleHTTP()
  return http.ListenAndServe(util.Address("", port), nil)
}

// Fetch the initial registry from the address
func connect(registry *Registry, address string) error {
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		return err
	}

	var reply *NodeList
	err = client.Call("Registry.Query", "", &reply)
	if err != nil {
		return err
	}

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
	gob.Register(registry.Self)
	return registry.Announce(message, &reply)
}

func Start(name string, hostname string, seeds []string, port int) (chan []byte, error) {
	address := util.Address(hostname, port)
	registry = NewRegistry(name)
	registry.Self = &Node{Name: name, Address: address}

  go server(registry, port)
  err := client(registry, seeds, port)
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
		Args:          data,
	}
	gob.Register(registry.Self)
	registry.Announce(message, &reply)
}
