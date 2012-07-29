package node_list

import (
  . "gossip/node"
)

type NodeList struct {
	Self  *Node
	Nodes map[string]*Node
}
