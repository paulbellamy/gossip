package node

import (
  "fmt"
)

type Node struct {
  Name string
  Address string
  Services map[string]int
}

func (n *Node) String() string {
  return fmt.Sprintf(" %s: %v", n.Address, n.Services)
}
