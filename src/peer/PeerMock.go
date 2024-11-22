package peer

import (
	"bittorrent/common"
	"fmt"
	"net"
	"sync"
)

type PeerMock struct {
	Address common.Address
}

func (peer PeerMock) Torrent(waitGroup *sync.WaitGroup) error {
	if waitGroup != nil {
		defer waitGroup.Done()
	}

	address := peer.Address.Ip + ":" + peer.Address.Port
	listener, err := net.Listen("tcp", address)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("MOCK " + address + ": Mock client listening for new connections")
	for {
		connection, err := listener.Accept()

		if err != nil {
			return err
		}

		go HandleConnection(connection)
	}
}

func (peer PeerMock) ActiveTorrent(wg *sync.WaitGroup, address common.Address) error {
	if wg != nil {
		defer wg.Done()
	}

	connection, err := net.Dial("tcp", address.Ip+":"+address.Port)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	message := "america:" + connection.LocalAddr().String()
	connection.Write([]byte(message))

	go HandleConnection(connection)

	return nil
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("MOCK " + conn.LocalAddr().String() + ": Connection established with " + conn.RemoteAddr().String())
	for {
		buffer := make([]byte, 1024)
		length, err := conn.Read(buffer)

		if err != nil {
			return
		}

		if length == 0 {
			continue
		}

		message := buffer[:length]

		fmt.Println("MOCK " + conn.LocalAddr().String() + ": Message received -> " + string(message))
	}
}
