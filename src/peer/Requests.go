package peer

import (
	"bittorrent/common"
	"bittorrent/tracker"
	"fmt"
	"net"
	"strings"
	"time"
)

func requestTracker(notificationChannel chan interface{}, tracker tracker.Tracker, request tracker.TrackRequest, timeToWait int) {
	time.Sleep(time.Second * time.Duration(timeToWait)) // Wait for the specified time
	response, err := tracker.Track(request)

	// Handle tracker's error
	if err != nil {
		fmt.Println(err.Error())
		notificationChannel <- trackerResponseNotification{Successful: false}
		return
	}

	notificationChannel <- trackerResponseNotification{Response: response, Successful: true}
}

func requestDownload(notificationChannel chan interface{}, timeToWait int) {

	time.Sleep(time.Second * time.Duration(timeToWait)) // Wait for the specified time

	// Here goes the logic involving downloading
	// ...
	//

	notificationChannel <- downloadNotification{}
}

func requestPeerUp(notificationChannel chan interface{}, id string, address common.Address) {
	connection, err := net.Dial("tcp", address.Ip+":"+address.Port)

	// Check if connection could not be established, if so then stop
	if err != nil {
		return
	}

	err = startHandshake(connection)

	// Check if handshaking could not be done, if so then stop
	if err != nil {
		return
	}

	fmt.Println("PEER: Connection established from: " + connection.LocalAddr().String())

	notificationChannel <- peerUpNotification{
		Id:         id,
		Connection: connection,
	}
}

func requestPeerListen(notificationChannel chan interface{}, address common.Address) {
	listener, err := net.Listen("tcp", address.Ip+":"+address.Port)

	fmt.Println("PEER: Start listening")
	// TODO: Properly handle error here
	if err != nil {
		fmt.Println("Peer could not start listening")
		notificationChannel <- killNotification{}
	}

	for {
		connection, err := listener.Accept()

		if err != nil {
			continue
		}

		go receiveHandshake(notificationChannel, connection)
	}
}

func startHandshake(connection net.Conn) error {
	// TODO: Perform a bittorrent starter-handshake
	message := []byte("Hi Neighbor!")
	err := common.ReliableWrite(connection, message)
	return err
}

func receiveHandshake(notificationChannel chan interface{}, connection net.Conn) error {
	// TODO: Perform a bittorrent receiver-handshake
	for {
		buffer := make([]byte, 1024)
		bytesRead, err := connection.Read(buffer)

		if err != nil {
			return err
		}

		if bytesRead == 0 {
			continue
		}

		message := string(buffer[:bytesRead])
		messageList := strings.Split(message, ":")

		if len(messageList) < 3 {
			continue
		}

		id := messageList[0]

		fmt.Println("PEER: Received a handshake from " + connection.RemoteAddr().String())

		notificationChannel <- peerUpNotification{
			Id:         id,
			Connection: connection,
		}

		return nil
	}
}
