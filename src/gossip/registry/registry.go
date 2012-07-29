package registry

import (
  "log"
  . "gossip/node"
  . "gossip/message"
  "net/rpc"
)

type Registry struct {
  Self *Node
  Nodes map[string]*Node
  Seen map[string] map[uint64]bool
}

func NewRegistry(name string) *Registry {
  reg := &Registry{}
  reg.Nodes = map[string]*Node{}
  reg.Seen = map[string] map[uint64]bool{}
  return reg
}

// For other nodes to fetch a copy of this registry
// the 'regex' argument is not really used yet.
func (t *Registry) Query(regex string, reply *Registry) error {
  log.Println("Query:",*t)

  *reply = *t
  return nil
}

// Announce a message to the cluster
func (t *Registry) Announce(message *Message, reply *int) error {
  log.Println("Announce:", message.Origin.Name)

  *reply = 0

  if _, exists := t.Seen[message.Origin.Name]; !exists {
    // Never seen a message from this node before
    t.Seen[message.Origin.Name] = map[uint64]bool{}
  }

  if _, exists := t.Seen[message.Origin.Name][message.Seq]; !exists {
    // We haven't seen this message before
    *reply = 1
    t.Seen[message.Origin.Name][message.Seq] = true

    // Add the new node to our registry
    // This should be determined by the actual message
    if message.ServiceMethod == "Registry.AddNode" {
      AddNode(t, message.Args.(*Node))
    }

    // Announce the new message to two other nodes
    var sent = 0
    var ok int
    var err error
    for _, old_node := range t.Nodes {
      if old_node.Name != message.Origin.Name {
        ok, err = publish(message, old_node.Address)
        if err == nil && ok == 1 {
          log.Println("Announced",message.Seq,"from",message.Origin.Name,"to",old_node.Name)
          sent++
        } else if err != nil {
          log.Println("Error Announcing",message.Seq,"from",message.Origin.Name,"to",old_node.Name,",",err)
        }

        if sent >= 2 {
          break
        }
      }
    }
  }

  return nil
}

func publish(message *Message, address string) (int, error) {
  client, err := rpc.DialHTTP("tcp", address)
  if err != nil {
    return 0, err
  }

  var reply int
  err = client.Call("Registry.Announce", message, &reply)
  return reply, err
}

// Merge remote registry into local registry
func MergeRegistries(local, remote *Registry) {
  for _, node := range remote.Nodes {
    AddNode(local, node)
  }
  AddNode(local, remote.Self)
}

func AddNode(reg *Registry, node *Node) {
  log.Println("AddNode:",node.Name)
  reg.Nodes[node.Name] = node
}
