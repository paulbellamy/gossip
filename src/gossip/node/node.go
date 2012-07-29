package node

import (
  "fmt"
)

type Node struct {
  Name string
  Address string
}

func (n *Node) String() string {
  return fmt.Sprintf("%s", n.Address)
}
