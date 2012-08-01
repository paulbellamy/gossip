package registry

import (
	. "gossip/message"
	. "gossip/node"
	. "gossip/node_list"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Messager struct {
	Seen map[string]map[uint64]bool
	Data chan []byte
}

type Registry struct {
	*Messager
	*NodeList
}

func NewRegistry(name string) *Registry {
	reg := &Registry{
		&Messager{
			Seen: map[string]map[uint64]bool{},
			Data: make(chan []byte),
		},
		&NodeList{
			Self:  nil,
			Nodes: map[string]*Node{},
		},
	}
	return reg
}

// For other nodes to fetch a copy of this registry
// the 'regex' argument is not really used yet.
func (t *Registry) Query(regex string, reply *NodeList) error {
	*reply = NodeList{
		Self:  t.Self,
		Nodes: t.Nodes,
	}
	return nil
}

// Announce a message to the cluster
func (t *Registry) Announce(message *Message, reply *int) error {
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
		if message.Origin.Name != t.Self.Name {
			if message.ServiceMethod == "Registry.AddNode" {
				args := message.Args.(map[string]interface{})
				AddNode(t, &Node{
					Name:    args["Name"].(string),
					Address: args["Address"].(string),
				})
			} else if message.ServiceMethod == "Registry.Data" {
				Data(t, []byte(message.Args.(string)))
			}
		}

		// Announce the new message to two other nodes
		go t.forward(message)
	}

	return nil
}

// Announce the new message to two other nodes
func (t *Registry) forward(message *Message) {
	var sent = 0
	var ok int
	var err error
	for _, old_node := range t.Nodes {
		if old_node.Name != message.Origin.Name && old_node.Name != t.Self.Name {
			ok, err = publish(message, old_node.Address)
			if err == nil && ok == 1 {
				sent++
			}

			if sent >= 2 {
				break
			}
		}
	}
}

func publish(message *Message, address string) (int, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return 0, err
	}

	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	if err != nil {
		return 0, err
	}

	var reply int
	err = client.Call("Registry.Announce", message, &reply)

	client.Close()
	return reply, err
}

func AddNode(reg *Registry, node *Node) {
	reg.Nodes[node.Name] = node
}

func Data(reg *Registry, data []byte) {
	reg.Data <- data
}

// Merge remote registry into local registry
func MergeRegistries(local *Registry, remote *NodeList) {
	for _, node := range remote.Nodes {
		AddNode(local, node)
	}
	AddNode(local, remote.Self)
}
