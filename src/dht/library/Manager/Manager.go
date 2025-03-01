package Manager

import "bittorrent/dht/library/BruteChord/Core"

type IManagerRPC[T Core.Contact] interface {
	GetActiveNodesIds() []Core.ChordHash
	GetNodeStateRPC(nodeId Core.ChordHash) string
}
