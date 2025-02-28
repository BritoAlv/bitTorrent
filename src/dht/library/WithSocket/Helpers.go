package WithSocket

import (
	"bittorrent/dht/library/BruteChord/Core"
	"encoding/gob"
	"net"
	"strings"
)

func (s *SocketServerClient) listenUDP() {
	for {
		buffer := make([]byte, 1024)
		n, _, err := s.listenerUDP.ReadFromUDP(buffer)
		if err != nil {
			s.logger.WriteToFileError("Error reading from UDP %v", err)
		}
		var notification Core.Notification[SocketContact]
		gobDecoder := gob.NewDecoder(strings.NewReader(string(buffer[:n])))
		err = gobDecoder.Decode(&notification)
		if err != nil {
			s.logger.WriteToFileError("Error decoding the notification %v", err)
		}
		s.communicationChannel <- notification
	}
}

func (s *SocketServerClient) listenTCP() {
	for {
		conn, err := s.listenerTCP.Accept()
		if err != nil {
			s.logger.WriteToFileError("Accepting the connection %v", err)
		}
		go s.handleConnection(conn)
	}
}

func (s *SocketServerClient) handleConnection(conn net.Conn) {
	// Read the data from the connection, and convert it to a Notification somehow, and pass this notification through the channel.
	s.logger.WriteToFileOK("Handling connection from %v", conn.RemoteAddr())
	gobDecoder := gob.NewDecoder(conn)
	var notification Core.Notification[SocketContact]
	err := gobDecoder.Decode(&notification)
	if err != nil {
		s.logger.WriteToFileError("Error decoding the notification %v", err)
	}
	s.communicationChannel <- notification
	err = conn.Close()
	if err != nil {
		s.logger.WriteToFileError("Error closing the connection %v", err)
	}
}

func createBroadcastAddress(ip string) string {
	parts := strings.Split(ip, ".")
	parts[len(parts)-1] = "255"
	return strings.Join(parts, ".")
}
