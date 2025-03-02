package Manager

import (
	"bittorrent/dht/library/BruteChord/Core"
	"bittorrent/dht/library/MonitorHand"
	"bittorrent/dht/library/WithSocket"
	"encoding/gob"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type HttpManager struct {
	Ports      []string
	KnownNodes Core.SafeMap[Core.ChordHash, string]
	Monitor    MonitorHand.MonitorHand[WithSocket.SocketContact]
}

func NewHttpManager(ports []string) *HttpManager {
	var httpManager HttpManager
	httpManager.Ports = ports
	httpManager.KnownNodes = *Core.NewSafeMap(make(map[Core.ChordHash]string))
	httpManager.Monitor = *MonitorHand.NewMonitorHand[WithSocket.SocketContact]("HttpManagerMonitor")
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
			h.KnownNodes.Set(Core.ChordHash(nodeId), state)
			h.Monitor.AddContact(WithSocket.SocketContact{
				NodeId: Core.ChordHash(nodeId),
				Addr:   nil,
			}, time.Now())
		}
		time.Sleep(1 * time.Second)
		for _, nodeId := range h.KnownNodes.GetKeys() {
			contact := WithSocket.SocketContact{
				NodeId: nodeId,
				Addr:   nil,
			}
			if !h.Monitor.CheckAlive(contact, 5) {
				h.KnownNodes.Delete(nodeId)
				h.Monitor.DeleteContact(contact)
			}
		}
	}
}

func (h *HttpManager) GetActiveNodesIds() []Core.ChordHash {
	return h.KnownNodes.GetKeys()
}

func (h *HttpManager) GetNodeStateRPC(nodeId Core.ChordHash) string {
	state, _ := h.KnownNodes.Get(nodeId)
	return state
}
