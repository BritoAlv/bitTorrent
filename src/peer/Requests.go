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

	// TODO: Properly handle tracker's error
	// Handle tracker's error
	if err != nil {
		fmt.Println("An error occurred contacting the tracker")
	}

	notificationChannel <- trackerResponseNotification{Response: response}
}

func requestDownload(notificationChannel chan interface{}, timeToWait int) {

	time.Sleep(time.Second * time.Duration(timeToWait)) // Wait for the specified time

	// Here goes the logic involving downloading
	// ...
	//

	notificationChannel <- downloadNotification{}
}

func requestPeerUp(notificationChannel chan interface{}, address common.Address, id string) {
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
		Address:    address,
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
	err := reliableWrite(connection, message)
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
		ip := messageList[1]
		port := messageList[2]

		fmt.Println("PEER: Received a handshake from " + connection.RemoteAddr().String())

		notificationChannel <- peerUpNotification{
			Address: common.Address{
				Ip:   ip,
				Port: port,
			},
			Id:         id,
			Connection: connection,
		}

		return nil
	}
}

func reliableWrite(connection net.Conn, message []byte) error {
	totalWritten := 0

	for totalWritten < len(message) {
		bytesWritten, err := connection.Write(message[totalWritten:])

		if err != nil {
			return err
		}

		totalWritten += bytesWritten
	}

	return nil
}
