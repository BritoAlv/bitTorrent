package TrackerNode

import (
	"bittorrent/common"
	"bittorrent/dht/library/BruteChord/Core"
	"bytes"
	"encoding/gob"
)

func (tracker *HttpTracker) InfoHashToChordKey(infoHash [20]byte) Core.ChordHash {
	sum := 0
	for i := 0; i < 20; i++ {
		sum += int(infoHash[i])
	}
	return Core.ChordHash(sum % (1 << Core.NumberBits))
}

func (tracker *HttpTracker) EncodePeerList(peers map[string]common.Address) []byte {
	var buf bytes.Buffer
	gob.NewEncoder(&buf)
	return buf.Bytes()
}

func (tracker *HttpTracker) DecodePeerList(data []byte) map[string]common.Address {
	var peers map[string]common.Address
	buf := bytes.NewBuffer(data)
	gob.NewDecoder(buf).Decode(&peers)
	return peers
}
