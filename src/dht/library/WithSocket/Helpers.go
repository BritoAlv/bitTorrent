package WithSocket

import (
	"bittorrent/dht/library/BruteChord/Core"
	"encoding/binary"
	"encoding/gob"
	"net"
)

func GetIpFromInterface(networkInterface string) (string, string) {
	itf, _ := net.InterfaceByName(networkInterface) //here your interface
	item, _ := itf.Addrs()
	for _, addr := range item {
		switch v := addr.(type) {
		case *net.IPNet:
			if !v.IP.IsLoopback() {
				if v.IP.To4() != nil { //Verify if IP is IPV4
					ip := v.IP
					if v.IP.To4() == nil {
						return "", ""
					}
					broadIP := make(net.IP, len(v.IP.To4()))
					binary.BigEndian.PutUint32(broadIP, binary.BigEndian.Uint32(v.IP.To4())|^binary.BigEndian.Uint32(net.IP(v.Mask).To4()))
					return ip.String(), broadIP.String()
				}
			}
		}
	}
	return "", ""
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
		panic(err)
	}
	s.communicationChannel <- notification
	err = conn.Close()
	if err != nil {
		s.logger.WriteToFileError("Error closing the connection %v", err)
	}
}
