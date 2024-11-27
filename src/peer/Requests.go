package peer

import (
	"bittorrent/common"
	"bittorrent/tracker"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const _REQUEST_MESSAGE = 6

const _HANDSHAKE_LENGTH = 30

func performTrack(notificationChannel chan interface{}, tracker tracker.Tracker, request tracker.TrackRequest, timeToWait int) {
	time.Sleep(time.Second * time.Duration(timeToWait)) // Wait for the specified time
	response, err := tracker.Track(request)

	// Handle tracker's error
	if err != nil {
		fmt.Println(err.Error())
		notificationChannel <- trackNotification{Successful: false}
		return
	}

	notificationChannel <- trackNotification{Response: response, Successful: true}
}

func performDownload(notificationChannel chan interface{}, timeToWait int) {
	time.Sleep(time.Second * time.Duration(timeToWait)) // Wait for the specified time
	notificationChannel <- downloadNotification{}
}

func performAddPeer(notificationChannel chan interface{}, sourceId string, targetId string, address common.Address) {
	connection, err := net.Dial("tcp", address.Ip+":"+address.Port)

	// Check if connection could not be established, if so then stop
	if err != nil {
		return
	}

	go performReadFromPeer(notificationChannel, connection, true, sourceId, targetId)

	err = sendHandshake(connection, sourceId)
	// Check if handshaking could not be done, if so then stop
	if err != nil {
		connection.Close()
		return
	}
}

func performListen(notificationChannel chan interface{}, address common.Address, sourceId string) {
	listener, err := net.Listen("tcp", address.Ip+":"+address.Port)

	fmt.Println("PEER: Start listening")
	// TODO: Properly handle error here
	if err != nil {
		fmt.Println("Peer could not start listening: " + err.Error())
		notificationChannel <- killNotification{}
		return
	}

	for {
		connection, err := listener.Accept()

		if err != nil {
			continue
		}

		go performReadFromPeer(notificationChannel, connection, false, sourceId, "")
	}
}

func performReadFromPeer(notificationChannel chan interface{}, connection net.Conn, active bool, sourceId string, targetId string) {
	wasHandshakeMade := false

	// TODO: Time out connections
	for {
		if !wasHandshakeMade {
			//TODO: Properly perform a handshake. Notice that peer-id must extracted/checked
			bytes, err := common.ReliableRead(connection, _HANDSHAKE_LENGTH)
			if err != nil {
				fmt.Println("PEER: an error occurred while reading from neighbor: " + err.Error())
				return
			}

			fmt.Println("PEER: Receive handshake message: " + string(bytes))
			if string(bytes[:10]) == "VivaaCubaa" {
				receivedId := string(bytes[10:_HANDSHAKE_LENGTH])

				if !active {
					targetId = receivedId
					sendHandshake(connection, sourceId)
				} else if active && targetId != receivedId {
					fmt.Println("PEER: invalid id")
					return
				}

				notificationChannel <- addPeerNotification{PeerId: receivedId, Connection: connection}
				wasHandshakeMade = true
			} else {
				fmt.Println("PEER: invalid handshake message")
				return
			}
		} else {
			bytes, err := common.ReliableRead(connection, 1)
			// TODO: Handle error (a good approach is to send a removePeerNotification)
			if err != nil {
				fmt.Println("PEER: an error occurred while reading from neighbor: " + err.Error())
				return
			}

			messageLength := int(bytes[0])

			bytes, err = common.ReliableRead(connection, messageLength)
			if err != nil {
				fmt.Println("PEER: an error occurred while reading from neighbor: " + err.Error())
				return
			}

			notification := handlePeerMessage(bytes, targetId)
			notificationChannel <- notification
		}
	}
}

func performDownloadFromPeer(notificationChannel chan interface{}, connection net.Conn, index int, offset int, length int) {
	// Message format: message-length;message-type;index;offset;length
	// TODO: Send the message in the official with big-endian and all that stuff
	message := strconv.Itoa(_REQUEST_MESSAGE) + ";" + strconv.Itoa(index) + ";" + strconv.Itoa(offset) + ";" + strconv.Itoa(length) + ";"
	bytes := []byte(message)
	lengthPrefix := byte(len(bytes))
	bytes = append([]byte{lengthPrefix}, bytes...)
	err := common.ReliableWrite(connection, bytes)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func sendHandshake(connection net.Conn, sourceId string) error {
	// TODO: Perform a bittorrent starter-handshake
	message := []byte("VivaaCubaa" + sourceId)
	err := common.ReliableWrite(connection, message)
	return err
}

// Receives a message and returns a notification
func handlePeerMessage(message []byte, peerId string) interface{} {
	// TODO correctly parse message's payload
	// messageType := message[0]
	// payload := message[1:]
	splits := strings.Split(string(message), ";")

	fmt.Printf("PEER: Handling a peer message. Payload: %v \n", splits[1:])

	switch splits[0] {
	case "0":
		fmt.Println("PEER: Receive a choke message")
		return peerChokeNotification{
			PeerId: peerId,
			Choke:  true,
		}
	case "1":
		fmt.Println("PEER: Receive a unchoke message")
		return peerChokeNotification{
			PeerId: peerId,
			Choke:  false,
		}
	case "2":
		fmt.Println("PEER: Receive a interested message")
		return peerInterestedNotification{
			PeerId:     peerId,
			Interested: true,
		}
	case "3":
		fmt.Println("PEER: Receive a not interested message")
		return peerInterestedNotification{
			PeerId:     peerId,
			Interested: false,
		}
	case "4":
		fmt.Println("PEER: Receive a have message")
		return peerHaveNotification{
			PeerId: peerId,
			Index:  0,
		}
	case "5":
		fmt.Println("PEER: Receive a bitfield message")
		return peerBitfieldNotification{
			PeerId:   peerId,
			Bitfield: []bool{},
		}
	case "6":
		fmt.Println("PEER: Receive a request message")
		return peerRequestNotification{
			PeerId: peerId,
			Index:  0,
			Offset: 0,
			Length: 0,
		}
	case "7":
		fmt.Println("PEER: Receive a piece message")
		return peerPieceNotification{
			PeerId: peerId,
			Index:  0,
			Offset: 0,
			Bytes:  []byte{},
		}
	case "8":
		fmt.Println("PEER: Receive a cancel message")
		return peerCancelNotification{
			PeerId: peerId,
			Index:  0,
			Offset: 0,
			Length: 0,
		}
	default:
		fmt.Println("PEER: Error! Invalid message type")
		return killNotification{}
	}
}
