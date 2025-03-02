package Manager

import (
	"bittorrent/dht/library/BruteChord/Core"
	"encoding/gob"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type HttpManager struct {
	Ports      []string
	KnownNodes Core.SafeMap[Core.ChordHash, string]
}

func NewHttpManager(ports []string) *HttpManager {
	var httpManager HttpManager
	httpManager.Ports = ports
	httpManager.KnownNodes = *Core.NewSafeMap(make(map[Core.ChordHash]string))
	go httpManager.updateStates()
	return &httpManager
}

func (h *HttpManager) updateStates() {
	for {
		for _, port := range h.Ports {
			resp, err := http.Get("http://localhost:" + port + "/nodeState")
			if err != nil {
				continue
			}
			var state string
			decoder := gob.NewDecoder(resp.Body)
			err = decoder.Decode(&state)
			if err != nil {
				continue
			}
			nodeId, _ := strconv.Atoi(strings.Split(state, "\n")[0][6:])
			fmt.Println(nodeId)
			h.KnownNodes.Set(Core.ChordHash(nodeId), state)
		}
		time.Sleep(1 * time.Second)
	}
}

func (h *HttpManager) GetActiveNodesIds() []Core.ChordHash {
	return h.KnownNodes.GetKeys()
}

func (h *HttpManager) GetNodeStateRPC(nodeId Core.ChordHash) string {
	state, _ := h.KnownNodes.Get(nodeId)
	return state
}
