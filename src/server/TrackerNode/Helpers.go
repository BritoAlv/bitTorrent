package TrackerNode

import (
	"bittorrent/common"
	"bittorrent/dht/library/BruteChord/Core"
	"bytes"
	"encoding/gob"
	"fmt"
)

func (tracker *HttpTracker) InfoHashToChordKey(infoHash [20]byte) Core.ChordHash {
	sum := 0
	for i := 0; i < 20; i++ {
		sum += int(infoHash[i])
	}
	return Core.ChordHash(sum % (1 << Core.NumberBits))
}

func (tracker *HttpTracker) EncodePeerList(peers map[string]common.Address) []byte {
	fmt.Printf("Passed to encode this %v\n", peers)
	if peers == nil {
		panic("Passed Peers is nil")
	}
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(peers)
	if err != nil {
		panic(err)
	}
	bytesEncoded := buf.Bytes()
	fmt.Printf("Encoded this \n%v\n", bytesEncoded)
	return bytesEncoded
}

func (tracker *HttpTracker) DecodePeerList(data []byte) map[string]common.Address {
	var peers map[string]common.Address
	fmt.Printf("Received this to decode \n%v\n", data)
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&peers)
	if err != nil {
		panic(err)
	}
	return peers
}
