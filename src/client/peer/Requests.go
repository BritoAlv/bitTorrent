package peer

import (
	"bittorrent/client/messenger"
	"bittorrent/client/tracker"
	"bittorrent/common"
	"fmt"
	"net"
	"time"
)

func performTrack(notificationChannel chan interface{}, tracker tracker.Tracker, request common.TrackRequest, timeToWait int) {
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

func performAddPeer(notificationChannel chan interface{}, sourceId string, targetId string, address common.Address, infohash [20]byte) {
	connection, err := net.Dial("tcp", address.Ip+":"+address.Port)

	// Check if connection could not be established, if so then stop
	if err != nil {
		return
	}

	go performReadFromPeer(notificationChannel, connection, true, sourceId, targetId, infohash)

	_messenger := messenger.New()
	err = _messenger.Write(connection, messenger.HandshakeMessage{
		Infohash: infohash,
		Id:       sourceId,
	})
	// Check if handshaking could not be done, if so then stop
	if err != nil {
		connection.Close()
		return
	}
}

func performListen(notificationChannel chan interface{}, address common.Address, sourceId string, infohash [20]byte) {
	listener, err := net.Listen("tcp", address.Ip+":"+address.Port)

	if err != nil {
		fmt.Println("ERROR: could not start listening: " + err.Error())
		notificationChannel <- killNotification{}
		return
	}

	fmt.Println("LOG: start listening")
	for {
		connection, err := listener.Accept()

		if err != nil {
			continue
		}

		go performReadFromPeer(notificationChannel, connection, false, sourceId, "", infohash)
	}
}

func performReadFromPeer(notificationChannel chan interface{}, connection net.Conn, active bool, sourceId string, targetId string, infohash [20]byte) {
	_messenger := messenger.New()
	wasHandshakeMade := false

	// TODO: Time out connections
	for {
		message, err := _messenger.Read(connection)
		if err != nil {
			fmt.Println("ERROR: an error occurred while reading from neighbor: " + err.Error())
			notificationChannel <- removePeerNotification{PeerId: targetId}
			return
		}

		switch castedMessage := message.(type) {
		case messenger.HandshakeMessage:
			if !wasHandshakeMade {
				if active && (targetId != castedMessage.Id || infohash != castedMessage.Infohash) {
					fmt.Println("ERROR: not expected id or infohash")
					notificationChannel <- removePeerNotification{PeerId: targetId}
					return
				}

				if !active {
					targetId = castedMessage.Id
					err := _messenger.Write(connection, messenger.HandshakeMessage{
						Infohash: infohash,
						Id:       sourceId,
					})
					if err != nil {
						fmt.Println("ERROR: an error occurred while reading from neighbor: " + err.Error())
						notificationChannel <- removePeerNotification{PeerId: targetId}
						return
					}
				}
				wasHandshakeMade = true

				fmt.Println("LOG: handshake performed with: " + targetId)
				// Notify to add a new peer
				notificationChannel <- addPeerNotification{
					PeerId:     targetId,
					Connection: connection,
				}

				// Send bitfields after handshake correctly performed
				notificationChannel <- sendBitfieldNotification{
					PeerId: targetId,
				}
			} else {
				fmt.Println("ERROR: handshake was already done")
				notificationChannel <- removePeerNotification{PeerId: targetId}
				return
			}
		case messenger.ChokeMessage:
			notificationChannel <- peerChokeNotification{
				PeerId: targetId,
				Choke:  true,
			}
		case messenger.UnchokeMessage:
			notificationChannel <- peerChokeNotification{
				PeerId: targetId,
				Choke:  false,
			}
		case messenger.InterestedMessage:
			notificationChannel <- peerInterestedNotification{
				PeerId:     targetId,
				Interested: true,
			}
		case messenger.NotInterestedMessage:
			notificationChannel <- peerInterestedNotification{
				PeerId:     targetId,
				Interested: false,
			}
		case messenger.HaveMessage:
			notificationChannel <- peerHaveNotification{
				PeerId: targetId,
				Index:  castedMessage.Index,
			}
		case messenger.BitfieldMessage:
			notificationChannel <- peerBitfieldNotification{
				PeerId:   targetId,
				Bitfield: castedMessage.Bitfield,
			}
		case messenger.RequestMessage:
			notificationChannel <- peerRequestNotification{
				PeerId: targetId,
				Index:  castedMessage.Index,
				Offset: castedMessage.Offset,
				Length: castedMessage.Length,
			}
		case messenger.PieceMessage:
			notificationChannel <- peerPieceNotification{
				PeerId: targetId,
				Index:  castedMessage.Index,
				Offset: castedMessage.Offset,
				Bytes:  castedMessage.Bytes,
			}
		case messenger.CancelMessage:
			notificationChannel <- peerCancelNotification{
				PeerId: targetId,
				Index:  castedMessage.Index,
				Offset: castedMessage.Offset,
				Length: castedMessage.Length,
			}
		default:
			fmt.Println("ERROR: invalid message type")
			notificationChannel <- removePeerNotification{PeerId: targetId}
			return
		}
	}
}

func performSendRequestToPeer(notificationChannel chan interface{}, connection net.Conn, peerId string, index int, offset int, length int) {
	_messenger := messenger.New()

	err := _messenger.Write(connection, messenger.RequestMessage{
		Index:  index,
		Offset: offset,
		Length: length,
	})

	if err != nil {
		fmt.Println("ERROR: an error occurred while sending a request to neighbor: " + err.Error())
		notificationChannel <- removePeerNotification{PeerId: peerId}
		return
	}

	fmt.Println("LOG: send a request-message to neighbor: " + peerId)
}

func performSendBitfieldToPeer(notificationChannel chan interface{}, connection net.Conn, peerId string, bitfield []bool) {
	_messenger := messenger.New()

	err := _messenger.Write(connection, messenger.BitfieldMessage{Bitfield: bitfield})

	if err != nil {
		fmt.Println(err.Error())
		notificationChannel <- removePeerNotification{PeerId: peerId}
		return
	}

	fmt.Printf("LOG: send a bitfield-message to neighbor: %v -> %v\n", peerId, bitfield)
}
