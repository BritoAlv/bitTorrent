package WithSocket

import (
	"bittorrent/dht/library/BruteChord/Core"
	"reflect"
	"time"
)

type ManagerSocket struct {
	socket     *SocketServerClient
	nodeStates Core.SafeMap[Core.ChordHash, string]
	channel    chan Core.Notification[SocketContact]
}

func NewManagerSocket() *ManagerSocket {
	var managerSocket ManagerSocket
	managerSocket.channel = make(chan Core.Notification[SocketContact])
	managerSocket.nodeStates = *Core.NewSafeMap[Core.ChordHash, string](make(map[Core.ChordHash]string))
	managerSocket.socket = NewSocketServerClient(-1)
	managerSocket.socket.SetData(managerSocket.channel, -1)
	go managerSocket.listenStates()
	go managerSocket.updateStates()
	return &managerSocket
}

func (m *ManagerSocket) listenStates() {
	for {
		select {
		case notification := <-m.channel:
			if reflect.TypeOf(notification) == reflect.TypeOf(&Core.TellMeYourStateResponse[SocketContact]{}) {
				response := notification.(*Core.TellMeYourStateResponse[SocketContact])
				m.nodeStates.Set(response.Sender.GetNodeId(), response.State)
			}
		}
	}
}

func (m *ManagerSocket) updateStates() {
	for {
		time.Sleep(1 * time.Second)
		m.socket.SendRequestEveryone(Core.TellMeYourState[SocketContact]{
			QueryHost: m.socket.GetContact(),
		})
	}

}

func (m *ManagerSocket) GetActiveNodesIds() []Core.ChordHash {
	result := m.socket.announcer.activeKnown.GetKeys()
	filtered := make([]Core.ChordHash, 0)
	for _, id := range result {
		if id >= 0 {
			filtered = append(filtered, id)
		}
	}
	return filtered
}

func (m *ManagerSocket) GetNodeStateRPC(nodeId Core.ChordHash) string {
	state, ok := m.nodeStates.Get(nodeId)
	if !ok {
		return ""
	}
	return state
}
