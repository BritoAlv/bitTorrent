package peer

import (
	"bittorrent/client/messenger"
	"bittorrent/client/tracker"
	"bittorrent/common"
	"fmt"
	"net"
	"slices"
	"time"
)

const _REQUEST_MESSAGE = 6

const _HANDSHAKE_LENGTH = 30

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

func performAddPeer(notificationChannel chan interface{}, sourceId string, targetId string, address common.Address, infohash []byte) {
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

func performListen(notificationChannel chan interface{}, address common.Address, sourceId string, infohash []byte) {
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

		go performReadFromPeer(notificationChannel, connection, false, sourceId, "", infohash)
	}
}

func performReadFromPeer(notificationChannel chan interface{}, connection net.Conn, active bool, sourceId string, targetId string, infohash []byte) {
	_messenger := messenger.New()
	wasHandshakeMade := false

	// TODO: Time out connections
	for {
		// TODO: Should send removePeerNotification if an error occurs
		message, err := _messenger.Read(connection)
		if err != nil {
			fmt.Println("PEER: An error occurred while reading from neighbor: " + err.Error())
			return
		}

		switch castedMessage := message.(type) {
		case messenger.HandshakeMessage:
			if !wasHandshakeMade {
				if active && (targetId != castedMessage.Id || slices.Compare(infohash, castedMessage.Infohash) != 0) {
					fmt.Println("PEER: Not expected id or not expected infohash")
					return
				}

				if !active {
					targetId = castedMessage.Id
					err := _messenger.Write(connection, messenger.HandshakeMessage{
						Infohash: infohash,
						Id:       sourceId,
					})
					if err != nil {
						fmt.Println("PEER: An error occurred while reading from neighbor: " + err.Error())
						return
					}
				}
				wasHandshakeMade = true

				fmt.Println("PEER: Handshake performed with: " + targetId)
				// Notify to add a new peer
				notificationChannel <- addPeerNotification{
					PeerId:     targetId,
					Connection: connection,
				}
			} else {
				fmt.Println("PEER: Handshake was already done")
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
		}
	}
}

func performDownloadFromPeer(notificationChannel chan interface{}, connection net.Conn, index int, offset int, length int) {
	_messenger := messenger.New()
	err := _messenger.Write(connection, messenger.RequestMessage{
		Index:  index,
		Offset: offset,
		Length: length,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
