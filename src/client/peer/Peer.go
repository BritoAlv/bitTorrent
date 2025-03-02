package peer

import (
	"bittorrent/client/pieceManager"
	"bittorrent/client/tracker"
	"bittorrent/common"
	"bittorrent/fileManager"
	"bittorrent/torrent"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

// **Peer's structure**
type Peer struct {
	Id string

	//** Private properties
	NotificationChannel chan interface{}
	address             common.Address            // Peer's address
	listener            net.Listener              // Peer's listener
	torrentData         torrent.Torrent           // Torrent associated data
	peers               map[string]PeerInfo       // Neighbor peers. It's a <PeerId, PeerInfo> dictionary
	tempPeers           map[string]common.Address // Peers being currently processed, might or not be official neighbors. This property can be refactor in the future
	privateKey          *rsa.PrivateKey
	requestedChunks     map[[3]string]int

	// Interfaces
	tracker      tracker.Tracker
	fileManager  fileManager.FileManager
	pieceManager pieceManager.PieceManager

	downloaded        bool               // File downloaded flag
	getAbsoluteOffset func(int, int) int // Function that calculates the absolute offset from index and relative-offset
}

func New(id string, listener net.Listener, torrent torrent.Torrent, downloadDirectory string, encrypted bool) (Peer, error) {
	splitAddress := strings.Split(listener.Addr().String(), ":")

	peer := Peer{}
	peer.Id = id
	peer.address = common.Address{Ip: splitAddress[0], Port: splitAddress[1]}
	peer.listener = listener
	peer.torrentData = torrent
	peer.NotificationChannel = make(chan interface{}, 1000)
	peer.peers = make(map[string]PeerInfo)
	peer.tempPeers = make(map[string]common.Address)
	peer.requestedChunks = make(map[[3]string]int)

	peer.tracker = tracker.NewTracker(torrent.Announce)

	length := 0
	var files []common.FileInfo
	if torrent.Files == nil {
		files = []common.FileInfo{{
			Length: int(torrent.Length),
			Path:   "/" + torrent.Name,
		}}
		length = int(torrent.Length)
	} else {
		files = torrent.Files
		for _, info := range files {
			length += info.Length
		}
		downloadDirectory += "/" + torrent.Name
	}

	var err error
	peer.fileManager, err = fileManager.New(downloadDirectory, files)
	if err != nil {
		return Peer{}, err
	}

	peer.pieceManager = pieceManager.New(length, int(torrent.PieceLength), common.CHUNK_SIZE)

	peer.getAbsoluteOffset = func(index int, offset int) int {
		return index*int(torrent.PieceLength) + offset
	}

	err = peer.checkAllPieces()
	if err != nil {
		return Peer{}, err
	}

	peer.privateKey = nil
	if encrypted {
		peer.privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return Peer{}, errors.New("error generating private key")
		}
	}

	return peer, nil
}

// **Peer's methods**

func (peer *Peer) Torrent(externalWaitGroup *sync.WaitGroup) error {
	if externalWaitGroup != nil {
		defer externalWaitGroup.Done()
	}

	trackerRequest := common.TrackRequest{
		InfoHash: peer.torrentData.InfoHash,
		PeerId:   peer.Id,
		Ip:       peer.address.Ip,
		Port:     peer.address.Port,
		Left:     peer.bytesLeft(),
		// Event:    "started",
	}

	go performListen(peer.NotificationChannel, peer.listener, peer.Id, peer.torrentData.InfoHash, peer.privateKey)
	go performTrack(peer.NotificationChannel, peer.tracker, trackerRequest, 0)
	if !peer.downloaded {
		go performDownload(peer.NotificationChannel, 2)
	}

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		for message := range peer.NotificationChannel {
			switch notification := message.(type) {
			case KillNotification:
				peer.listener.Close()
				for _, peerStruct := range peer.peers {
					peerStruct.Connection.Close()
				}
				fmt.Println("LOG: Killed")
				return
			case trackNotification:
				peer.handleTrackResponseNotification(notification)
			case downloadNotification:
				peer.handleDownloadNotification()
			case writeNotification:
				peer.handleWriteNotification(notification)
			case pieceVerificationNotification:
				peer.handlePieceVerifiedNotification(notification)
			case addPeerNotification:
				peer.handleAddPeerNotification(notification)
			case removePeerNotification:
				peer.handleRemovePeerNotification(notification)
			case sendBitfieldNotification:
				peer.handleSendBitfieldNotification(notification)
			case peerRequestNotification:
				peer.handlePeerRequestNotification(notification)
			case peerBitfieldNotification:
				peer.handlePeerBitfieldNotification(notification)
			case peerPieceNotification:
				peer.handlePeerPieceNotification(notification)
			case peerHaveNotification:
				peer.handlePeerHaveNotification(notification)
			default:
				fmt.Println("ERROR: invalid notification-type")
			}
		}
	}()

	waitGroup.Wait()
	return nil
}

func (peer Peer) Status() (progress float32, peers int) {
	bitfield := peer.pieceManager.Bitfield()
	indexes := len(bitfield)
	downloaded := 0
	for _, value := range bitfield {
		if value {
			downloaded++
		}
	}
	return float32(downloaded) / float32(indexes), len(peer.peers)
}

func (peer *Peer) handleTrackResponseNotification(notification trackNotification) {
	// TODO: Properly handle the case when the notification was not successful
	fmt.Println("LOG: handling tracker response notification")
	const PEERS_LOWER_BOUND int = 20

	// Check if the file is not downloaded and neighbor-peers are less than the established bound
	if !peer.downloaded && len(peer.peers) < PEERS_LOWER_BOUND {
		for id, address := range notification.Response.Peers {
			if peer.isValidNeighbor(id, address) {
				go performAddPeer(peer.NotificationChannel, peer.Id, id, address, peer.torrentData.InfoHash, peer.privateKey)
			}
		}
	}

	left := peer.bytesLeft()

	go performTrack(peer.NotificationChannel, peer.tracker, common.TrackRequest{
		InfoHash: peer.torrentData.InfoHash,
		PeerId:   peer.Id,
		Ip:       peer.address.Ip,
		Port:     peer.address.Port,
		Left:     left,
		// Event:    "started",
	}, notification.Response.Interval)
}

func (peer *Peer) handleDownloadNotification() {
	// Constants
	const UNCHECKED_CHUNKS_PER_PIECE = 3
	const INDEX = 0
	const OFFSET = 1
	const LENGTH = 2

	fmt.Println("LOG: handling download notification")
	missing_pieces := peer.pieceManager.GetUncheckedPieces()

	if len(missing_pieces) == 0 {
		peer.downloaded = true
		return
	}

	uncheckedChunks := [][3]int{}
	for _, index := range missing_pieces {
		uncheckedChunks = append(uncheckedChunks, peer.pieceManager.GetUncheckedChunks(index, UNCHECKED_CHUNKS_PER_PIECE)...)
	}

	for _, chunk := range uncheckedChunks {
		for peerId, peerInfo := range peer.peers {
			indexStr := strconv.Itoa(chunk[INDEX])
			offsetStr := strconv.Itoa(chunk[OFFSET])
			requestChunkId := [3]string{peerId, indexStr, offsetStr}

			requestedCount, previouslyRequested := peer.requestedChunks[requestChunkId]

			if requestedCount == 0 && previouslyRequested {
				fmt.Print()
			}

			if peerInfo.Bitfield[chunk[INDEX]] && (!previouslyRequested || requestedCount == 0) {
				peer.requestedChunks[requestChunkId] = 20

				go performSendRequestToPeer(peer.NotificationChannel, peerInfo.Connection, peerId, chunk[INDEX], chunk[OFFSET], chunk[LENGTH])
				break
			}

			if previouslyRequested {
				peer.requestedChunks[requestChunkId]--
			}
		}
	}

	if !peer.downloaded {
		go performDownload(peer.NotificationChannel, 5)
	}
}

func (peer *Peer) handleAddPeerNotification(notification addPeerNotification) {
	_, contains := peer.peers[notification.PeerId]

	if contains {
		return
	}

	peer.peers[notification.PeerId] = PeerInfo{
		Connection: notification.Connection,
		Bitfield:   make([]bool, len(peer.pieceManager.Bitfield())),
		IsChoker:   false,
		IsChoked:   false,
		PublicKey:  notification.PublicKey,
	}
}

func (peer *Peer) handleRemovePeerNotification(notification removePeerNotification) {
	delete(peer.tempPeers, notification.PeerId) // Make sure any temporal peer is removed
	info, contains := peer.peers[notification.PeerId]

	if !contains {
		return
	}

	err := info.Connection.Close()
	if err != nil {
		fmt.Println("ERROR: an error occurred while closing connection " + err.Error())
	}

	delete(peer.peers, notification.PeerId)
	for key := range peer.requestedChunks {
		if key[0] == notification.PeerId {
			delete(peer.requestedChunks, key)
		}
	}
	fmt.Println("LOG: remove neighbor: " + notification.PeerId)
}

func (peer *Peer) handleSendBitfieldNotification(notification sendBitfieldNotification) {
	info, contains := peer.peers[notification.PeerId]

	if !contains {
		return
	}

	performSendBitfieldToPeer(peer.NotificationChannel, info.Connection, notification.PeerId, peer.pieceManager.Bitfield())
}

func (peer *Peer) handlePeerRequestNotification(notification peerRequestNotification) {
	info, contains := peer.peers[notification.PeerId]

	if !contains {
		return
	}

	fmt.Println("LOG: a request-message was received from: " + notification.PeerId)
	if !peer.pieceManager.VerifyPiece(notification.Index) {
		peer.handleRemovePeerNotification(removePeerNotification{notification.PeerId})
		return
	}

	start := peer.getAbsoluteOffset(notification.Index, notification.Offset)
	go performSendPieceToPeer(peer.NotificationChannel, info.Connection, peer.fileManager, notification.PeerId, notification.Index, notification.Offset, notification.Length, start, peer.peers[notification.PeerId].PublicKey)
}

func (peer *Peer) handlePeerBitfieldNotification(notification peerBitfieldNotification) {
	info, contains := peer.peers[notification.PeerId]

	if !contains {
		return
	}

	// Add bitfield to peer's info
	info.Bitfield = notification.Bitfield
	peer.peers[notification.PeerId] = info
	fmt.Printf("LOG: a bitfield-message was received from: %v\n", notification.PeerId)
}

func (peer *Peer) handlePeerPieceNotification(notification peerPieceNotification) {
	_, contains := peer.peers[notification.PeerId]
	if !contains {
		return
	}

	if peer.pieceManager.VerifyPiece(notification.Index) || peer.pieceManager.VerifyChunk(notification.Index, notification.Offset) {
		return
	}

	absoluteOffset := peer.getAbsoluteOffset(notification.Index, notification.Offset)
	go performWrite(peer.NotificationChannel, peer.fileManager, notification.Index, notification.Offset, absoluteOffset, notification.Bytes)

	fmt.Println("LOG: a piece-message was received from: " + notification.PeerId)
}

func (peer *Peer) handleWriteNotification(notification writeNotification) {

	checkedPiece := peer.pieceManager.CheckChunk(notification.Index, notification.Offset)

	if checkedPiece {
		pieceAbsoluteOffset := notification.Index * int(peer.torrentData.PieceLength)

		go performVerifyPiece(peer.NotificationChannel, peer.fileManager, notification.Index, pieceAbsoluteOffset, int(peer.torrentData.PieceLength), peer.torrentData.Pieces)
	}
}

func (peer *Peer) handlePieceVerifiedNotification(notification pieceVerificationNotification) {
	if notification.Verified {
		fmt.Printf("LOG: piece %v was verified\n", notification.Index)

		// Send have-message to all neighbor peers
		for peerId, peerInfo := range peer.peers {
			go performSendHaveToPeer(peer.NotificationChannel, peerInfo.Connection, peerId, notification.Index)
		}
	} else {
		fmt.Printf("LOG: piece %v was corrupted\n", notification.Index)
		peer.pieceManager.UncheckPiece(notification.Index)
	}
}

func (peer *Peer) handlePeerHaveNotification(notification peerHaveNotification) {
	info, contains := peer.peers[notification.PeerId]

	if !contains {
		return
	}

	// Update peer's info's bitfield
	info.Bitfield[notification.Index] = true
	peer.peers[notification.PeerId] = info
	fmt.Printf("LOG: a have-message with index %v was received from: %v \n", notification.Index, notification.PeerId)
}

func (peer *Peer) isValidNeighbor(neighborId string, address common.Address) bool {
	// Check neighbor's id and address are different from peer's
	if peer.Id == neighborId || peer.address == address {
		return false
	}

	// Check neighbor is not already registered
	_, contains := peer.peers[neighborId]
	if contains {
		return false
	}

	// Check neighbor address is not already registered
	addressStr := address.Ip + ":" + address.Port
	for _, peerInfo := range peer.peers {
		if peerInfo.Connection.RemoteAddr().String() == addressStr {
			return false
		}
	}

	// Check the id/address are not being temporary processed
	for peerId, tempAddress := range peer.tempPeers {
		tempAddressStr := tempAddress.Ip + ":" + tempAddress.Port
		if peerId == neighborId || tempAddressStr == addressStr {
			return false
		}
	}

	// Add address to temporary processed ones
	peer.tempPeers[neighborId] = address
	return true
}

// Calculate the amount of bytes left to download
func (peer *Peer) bytesLeft() int {
	// Calculate the amount of bytes left to download
	left := 0
	if !peer.downloaded {
		for range peer.pieceManager.GetUncheckedPieces() {
			left += int(peer.torrentData.PieceLength)
		}
	}
	return left
}

func (peer *Peer) checkAllPieces() error {
	for index := range len(peer.torrentData.Pieces) / 20 {
		pieceAbsoluteOffset := peer.getAbsoluteOffset(index, 0)
		hash, err := getPieceHash(peer.fileManager, index, pieceAbsoluteOffset, int(peer.torrentData.PieceLength))
		if err != nil {
			return err
		}

		isHashValid := checkPieceHash(index, hash, peer.torrentData.Pieces)
		if !isHashValid {
			peer.pieceManager.UncheckPiece(index)
		}
	}

	if len(peer.pieceManager.GetUncheckedPieces()) == 0 {
		peer.downloaded = true
	}

	return nil
}
