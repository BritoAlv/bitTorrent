package WithSocket

import (
	"bittorrent/common"
	"bittorrent/dht/library/BruteChord/Core"
	"encoding/gob"
	"net"
	"strconv"
)

type SocketServerClient struct {
	listenerTCP          net.Listener
	announcer            Announcer
	communicationChannel chan Core.Notification[SocketContact]
	logger               common.Logger
	nodeId               Core.ChordHash
}

func NewSocketServerClient(nodeId Core.ChordHash) *SocketServerClient {
	var socketServerClient SocketServerClient
	ip, _ := GetIpFromInterface(networkInterface)
	listenerTCP, err := net.Listen("tcp", ip+":0")
	if err != nil {
		panic(err)
	}
	socketServerClient.nodeId = nodeId
	socketServerClient.listenerTCP = listenerTCP
	socketServerClient.announcer = *NewAnnouncer(socketServerClient.GetContact())
	socketServerClient.logger = *common.NewLogger("ServerClientSocket" + strconv.Itoa(int(nodeId)) + ".txt")
	go socketServerClient.listen()
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
			continue
		}
		encoder := gob.NewEncoder(conn)
		err = encoder.Encode(&task.Data)
		if err != nil {
			s.logger.WriteToFileError("Error encoding the data %v", err)
			continue
		}
		err = conn.Close()
	}
}

func (s *SocketServerClient) SendRequestEveryone(data Core.Notification[SocketContact]) {
	targets := s.announcer.GetContacts()
	s.SendRequest(Core.ClientTask[SocketContact]{Targets: targets, Data: data})
}

func (s *SocketServerClient) listen() {
	for {
		conn, err := s.listenerTCP.Accept()
		if err != nil {
			s.logger.WriteToFileError("Error while accepting connection %v", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *SocketServerClient) handleConn(conn net.Conn) {
	decoder := gob.NewDecoder(conn)
	var notification Core.Notification[SocketContact]
	err := decoder.Decode(&notification)
	if err != nil {
		s.logger.WriteToFileError("Error decoding the notification %v", err)
		return
	}
	s.communicationChannel <- notification
	err = conn.Close()
	if err != nil {
		s.logger.WriteToFileError("Error closing the connection %v", err)
	}
}
