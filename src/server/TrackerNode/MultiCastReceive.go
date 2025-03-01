package TrackerNode

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func receiveFromMulticast(ip string, port string) {
	// Multicast group address and Port
	const MulticastIp = "224.0.0.1" // Replace with the actual multicast address
	const MulticastPort = 10000     // Replace with the actual Port used by the proxy

	// Resolve the UDP address
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", MulticastIp, MulticastPort))
	if err != nil {
		log.Println("Failed to resolve UDP address:", err)
		os.Exit(1)
	}

	// Create the UDP socket
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Println("Failed to listen on multicast UDP:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Set the socket options
	if err := conn.SetReadBuffer(1024); err != nil {
		log.Println("Failed to set read buffer:", err)
		os.Exit(1)
	}

	log.Printf("Listening for multicast messages on %s:%d\n", MulticastIp, MulticastPort)

	// Receive/respond loop
	buf := make([]byte, 1024)
	for {
		n, src, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println("Failed to read from UDP:", err)
			continue
		}
		bytes := buf[:n]
		log.Printf("Received %d bytes from %s: %s\n", n, src, string(bytes))

		message := strings.Split(string(bytes), ";")
		if len(message) != 2 {
			log.Println("Invalid message")
			continue
		}

		clientIp, clientPort := message[0], message[1]
		go func() {
			connection, err := net.Dial("tcp", fmt.Sprintf("%v:%v", clientIp, clientPort))
			if err != nil {
				log.Println("Failed to connect to client:", err)
				return
			}
			defer connection.Close()

			// TODO: must do EnsureWrite here
			_, err = connection.Write(fmt.Appendf(nil, "%v;%v", ip, port))
			if err != nil {
				log.Println("Failed to send message to client:", err)
				return
			}
			log.Println("Sent message to client")
		}()
	}
}
