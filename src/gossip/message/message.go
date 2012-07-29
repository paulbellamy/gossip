package message

import (
	. "gossip/node"
)

type Message struct {
	Origin *Node  // originating node
	Seq    uint64 // sequence number chosen by client
  ServiceMethod string
  Args interface{}
}
