package Manager

import (
	"bittorrent/dht/library/BruteChord/Core"
	"bittorrent/dht/library/MonitorHand"
	"bittorrent/dht/library/WithSocket"
	"bittorrent/server/TrackerNode"
	"encoding/gob"
	"fmt"
	"net/http"
	"strconv"
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

func (h *HttpManager) StateToString(state Core.NodeState[WithSocket.SocketContact]) string {
	result := "Node: " + strconv.Itoa(int(state.NodeId)) + "\n"
	result += "Successor: " + strconv.Itoa(int(state.SuccessorId)) + "\n"
	result += "Successor Data Replicas Are: " + "\n"
	for _, key := range Core.SortKeys(state.SuccessorData) {
		result += strconv.Itoa(int(key)) + " -> " + fmt.Sprintf("%v", TrackerNode.DecodePeerList(state.SuccessorData[key])) + "\n"
	}
	result += "SuccessorSuccessor: " + strconv.Itoa(int(state.SuccessorSuccessorId)) + "\n"
	result += "SuccessorSuccessor Data Replica:" + "\n"
	for _, key := range Core.SortKeys(state.SuccessorSuccessorData) {
		result += strconv.Itoa(int(key)) + " -> " + fmt.Sprintf("%v", TrackerNode.DecodePeerList(state.SuccessorSuccessorData[key])) + "\n"
	}
	result += "Predecessor: " + strconv.Itoa(int(state.PredecessorId)) + "\n"
	result += "Data stored:\n"
	for _, key := range Core.SortKeys(state.OwnData) {
		result += strconv.Itoa(int(key)) + " -> " + fmt.Sprintf("%v", TrackerNode.DecodePeerList(state.OwnData[key])) + "\n"
	}
	return result
}

func (h *HttpManager) updateStates() {
	for {
		for _, port := range h.Ports {
			resp, err := http.Get("http://localhost:" + port + "/nodeState")
			if err != nil {
				continue
			}
			var state Core.NodeState[WithSocket.SocketContact]
			decoder := gob.NewDecoder(resp.Body)
			err = decoder.Decode(&state)
			if err != nil {
				continue
			}
			nodeId := state.NodeId
			h.KnownNodes.Set(nodeId, h.StateToString(state))
			h.Monitor.AddContact(WithSocket.SocketContact{
				NodeId: nodeId,
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
