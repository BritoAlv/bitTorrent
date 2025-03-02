package Manager

import "bittorrent/dht/library/BruteChord/Core"

type IManagerRPC interface {
	GetActiveNodesIds() []Core.ChordHash
	GetNodeStateRPC(nodeId Core.ChordHash) string
}
