package registry

import (
  "fmt"
  . "node/node"
  "node/util"
)

type Registry struct {
  Self *Node
  Nodes map[string]*Node
}

func NewRegistry(name string) *Registry {
  reg := &Registry{}
  reg.Nodes = map[string]*Node{}
  return reg
}

func (t *Registry) QueryAll(address string, reply *int) error {
  fmt.Println("QueryAll:",address)

  for _, old_node := range t.Nodes {
    util.Publish(old_node, address)
  }
  util.Publish(t.Self, address)

  *reply = 1
  return nil
}

func (t *Registry) AddNode(node *Node, reply *int) error {
  fmt.Println("AddNode:", node.Name)

  if _, exists := t.Nodes[node.Name]; !exists && t.Self.Name != node.Name {
    t.Nodes[node.Name] = node

    // Announce the new node to everyone else
    for _, old_node := range t.Nodes {
      if old_node.Name != node.Name {
        util.Publish(node, old_node.Address)
      }
    }
  }

  *reply = 1
  return nil
}

func (t *Registry) RemoveNode(name string, reply *int) error {
  delete(t.Nodes, name)

  fmt.Println("Registry:", t)

  *reply = 1
  return nil
}
