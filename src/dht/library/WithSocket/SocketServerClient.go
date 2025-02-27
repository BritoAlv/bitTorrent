package WithSocket

import (
	"bittorrent/common"
	"bittorrent/dht/library/BruteChord/Core"
	"encoding/gob"
	"net"
)

type SocketServerClient struct {
	listenerTCP          net.Listener
	listenerUDP          net.UDPConn
	communicationChannel chan Core.Notification[SocketContact]
	logger               common.Logger
	nodeId               Core.ChordHash
}

func NewSocketServerClient(ip string, name string) *SocketServerClient {
	var socketServerClient SocketServerClient
	listenerTCP, err := net.Listen("tcp", ip)
	if err != nil {
		panic(err)
	}
	socketServerClient.listenerTCP = listenerTCP

	listenerUDP, err := net.ListenUDP("udp", nil)
	if err != nil {
		panic(err)
	}
	socketServerClient.listenerUDP = *listenerUDP
	socketServerClient.logger = *common.NewLogger(name)
	return &socketServerClient
}

func (s *SocketServerClient) GetContact() SocketContact {
	return SocketContact{
		NodeId: s.nodeId,
		Addr:   s.listenerTCP.Addr(),
	}
}

func (s *SocketServerClient) SetData(channel chan Core.Notification[SocketContact], NodeId Core.ChordHash) {
	s.nodeId = NodeId
	s.communicationChannel = channel
}

func (s *SocketServerClient) SendRequest(task Core.ClientTask[SocketContact]) {
	for _, target := range task.Targets {
		conn, err := net.Dial("tcp", target.Addr.String())
		if err != nil {
			s.logger.WriteToFileError("Error dialing the connection %v", err)
		}
		encoder := gob.NewEncoder(conn)
		err = encoder.Encode(task.Data)
		if err != nil {
			s.logger.WriteToFileError("Error encoding the data %v", err)
		}
		err = conn.Close()
	}
}

func (s *SocketServerClient) SendRequestEveryOne(data Core.Notification[SocketContact]) {
	broadcastAddress := createBroadcastAddress(s.listenerTCP.Addr().String())
	conn, err := net.Dial("udp", broadcastAddress)
	if err != nil {
		s.logger.WriteToFileError("Error dialing the connection %v", err)
	}
	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(data)
	if err != nil {
		s.logger.WriteToFileError("Error encoding the data %v", err)
	}
	err = conn.Close()
}
