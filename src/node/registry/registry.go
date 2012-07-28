package registry

import (
  "fmt"
)

type Node struct {
  Name string
  Address string
  Services map[string]int
}

type Registry struct {
  Nodes map[string]*Node
}

func NewRegistry(name string) *Registry {
  reg := &Registry{}
  reg.Nodes = map[string]*Node{}
  reg.Nodes[name] = &Node{Name: name}
  return reg
}

func (t *Registry) AddNode(node *Node, reply *int) error {
  if _, exists := t.Nodes[node.Name]; !exists {
    t.Nodes[node.Name] = node
  }

  fmt.Println("Registry:", t)

  *reply = 1
  return nil
}

func (t *Registry) RemoveNode(name string, reply *int) error {
  delete(t.Nodes, name)

  fmt.Println("Registry:", t)

  *reply = 1
  return nil
}
