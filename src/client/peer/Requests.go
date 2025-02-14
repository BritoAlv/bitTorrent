package peer

import (
	"bittorrent/client/messenger"
	"bittorrent/client/tracker"
	"bittorrent/common"
	"bittorrent/fileManager"
	"crypto/rsa"
	"crypto/sha1"
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

func performAddPeer(notificationChannel chan interface{}, sourceId string, targetId string, address common.Address, infohash [20]byte, sourcePrivateKey *rsa.PrivateKey) {
	connection, err := net.Dial("tcp", address.Ip+":"+address.Port)

	// Check if connection could not be established, if so then stop
	if err != nil {
		fmt.Println("ERROR: connection could not be established: " + err.Error())
		notificationChannel <- removePeerNotification{PeerId: targetId}
		return
	}

	go performReadFromPeer(notificationChannel, connection, true, sourceId, targetId, infohash, sourcePrivateKey)

	_messenger := messenger.New(nil, nil)
	err = _messenger.Write(connection, messenger.HandshakeMessage{
		Infohash:  infohash,
		Id:        sourceId,
		PublicKey: &sourcePrivateKey.PublicKey,
	})

	// Check if handshaking could not be done, if so then stop
	if err != nil {
		fmt.Println("ERROR: an error occurred while performing handshaking: " + err.Error())
		connection.Close()
		notificationChannel <- removePeerNotification{PeerId: targetId}
		return
	}
}

func performListen(notificationChannel chan interface{}, listener net.Listener, sourceId string, infohash [20]byte, sourcePrivateKey *rsa.PrivateKey) {
	fmt.Println("LOG: start listening")
	for {
		connection, err := listener.Accept()

		if err != nil {
			continue
		}

		go performReadFromPeer(notificationChannel, connection, false, sourceId, "", infohash, sourcePrivateKey)
	}
}

func performReadFromPeer(notificationChannel chan interface{}, connection net.Conn, active bool, sourceId string, targetId string, infohash [20]byte, sourcePrivateKey *rsa.PrivateKey) {
	_messenger := messenger.New(sourcePrivateKey, nil)
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
						Infohash:  infohash,
						Id:        sourceId,
						PublicKey: &sourcePrivateKey.PublicKey,
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
					PublicKey:  castedMessage.PublicKey,
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
	_messenger := messenger.New(nil, nil)

	err := _messenger.Write(connection, messenger.RequestMessage{
		Index:  index,
		Offset: offset,
		Length: length,
	})

	if err != nil {
		fmt.Println("ERROR: an error occurred while sending a request-message to neighbor: " + err.Error())
		notificationChannel <- removePeerNotification{PeerId: peerId}
		return
	}

	fmt.Println("LOG: send a request-message to neighbor: " + peerId)
}

func performSendBitfieldToPeer(notificationChannel chan interface{}, connection net.Conn, peerId string, bitfield []bool) {
	_messenger := messenger.New(nil, nil)

	err := _messenger.Write(connection, messenger.BitfieldMessage{Bitfield: bitfield})

	if err != nil {
		fmt.Println("ERROR: an error occurred while sending a bitfield-message to neighbor: " + err.Error())
		notificationChannel <- removePeerNotification{PeerId: peerId}
		return
	}

	fmt.Printf("LOG: send a bitfield-message to neighbor: %v\n", peerId)
}

func performSendPieceToPeer(notificationChannel chan interface{}, connection net.Conn, _fileManager fileManager.FileManager, peerId string, index int, offset int, length int, absoluteOffset int, targetPublicKey *rsa.PublicKey) {
	_messenger := messenger.New(nil, targetPublicKey)

	bytes, err := _fileManager.Read(absoluteOffset, length)
	if err != nil {
		fmt.Println("ERROR: an error occurred while reading from file: " + err.Error())
		return
	}

	err = _messenger.Write(connection, messenger.PieceMessage{
		Index:  index,
		Offset: offset,
		Bytes:  bytes,
	})

	if err != nil {
		fmt.Println("ERROR: an error occurred while sending a piece-message to neighbor: " + peerId)
		notificationChannel <- removePeerNotification{PeerId: peerId}
		return
	}
}

func performSendHaveToPeer(notificationChannel chan interface{}, connection net.Conn, peerId string, index int) {
	_messenger := messenger.New(nil, nil)

	err := _messenger.Write(connection, messenger.HaveMessage{Index: index})

	if err != nil {
		fmt.Println("ERROR: an error occurred while sending a have-message to neighbor: " + err.Error())
		notificationChannel <- removePeerNotification{PeerId: peerId}
		return
	}

	fmt.Printf("LOG: send a have-message to neighbor: %v\n", peerId)
}

func performWrite(notificationChannel chan interface{}, _fileManager fileManager.FileManager, index int, offset int, absoluteOffset int, bytes []byte) {
	err := _fileManager.Write(absoluteOffset, &bytes)
	if err != nil {
		fmt.Println("ERROR: an error occurred while writing the file: " + err.Error())
		return
	}

	notificationChannel <- writeNotification{
		Index:  index,
		Offset: offset,
	}
}

func performVerifyPiece(notificationChannel chan interface{}, _fileManager fileManager.FileManager, index int, pieceAbsoluteOffset int, pieceLength int, pieces []byte) {
	hash, err := getPieceHash(_fileManager, index, pieceAbsoluteOffset, pieceLength)
	if err != nil {
		return
	}

	isValidHash := checkPieceHash(index, hash, pieces)
	if isValidHash {
		notificationChannel <- pieceVerificationNotification{Index: index, Verified: true}
	} else {
		notificationChannel <- pieceVerificationNotification{Index: index, Verified: false}
	}
}

func getPieceHash(_fileManager fileManager.FileManager, index int, pieceAbsoluteOffset int, pieceLength int) ([20]byte, error) {
	bytes, err := _fileManager.Read(pieceAbsoluteOffset, pieceLength)

	if err != nil {
		// Check if the reading attempt was outside of the file bounds, if so the expected bytes are not yet downloaded
		_, isOutsideOfFileBounds := err.(fileManager.OutsideOfFileBoundsError)

		if isOutsideOfFileBounds {
			bytes = []byte{}
		} else {
			return [20]byte{}, err
		}
	}

	return sha1.Sum(bytes), nil
}

func checkPieceHash(index int, hash [20]byte, pieces []byte) bool {
	hashIndex := index * 20
	pieceHash := pieces[hashIndex : hashIndex+20]
	return hash == [20]byte(pieceHash)
}
