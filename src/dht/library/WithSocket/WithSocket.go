package WithSocket

import (
	"bittorrent/dht/library/BruteChord/Core"
	"net"
)

type SocketContact struct {
	NodeId Core.ChordHash
	Addr   net.Addr
}

func NewSocketContact(nodeId Core.ChordHash, addr net.Addr) SocketContact {
	return SocketContact{NodeId: nodeId, Addr: addr}
}

func (s SocketContact) GetNodeId() Core.ChordHash {
	return s.NodeId
}