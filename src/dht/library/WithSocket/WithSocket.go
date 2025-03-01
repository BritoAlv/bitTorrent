package WithSocket

import (
	"bittorrent/dht/library/BruteChord/Core"
	"bittorrent/dht/library/MonitorHand"
	"net"
	"strconv"
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

func NewNodeSocket() *Core.BruteChord[SocketContact] {
	randomId := Core.GenerateRandomBinaryId()
	randomIdStr := strconv.Itoa(int(randomId))
	socketServerClient := NewSocketServerClient(randomId)
	monitorHand := MonitorHand.NewMonitorHand[SocketContact]("Monitor" + randomIdStr)
	nodeSocket := Core.NewBruteChord[SocketContact](socketServerClient, socketServerClient, monitorHand, randomId)
	return nodeSocket
}
