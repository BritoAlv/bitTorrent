package InMemory

import (
	"bittorrent/dht/library/BruteChord/Core"
)

type ContactInMemory struct {
	Id     string
	NodeId Core.ChordHash
}

func (i ContactInMemory) GetNodeId() Core.ChordHash {
	return i.NodeId
}
