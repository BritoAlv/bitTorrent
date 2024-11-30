package peer

import (
	"bittorrent/client/fileManager"
	"bittorrent/client/pieceManager"
	"bittorrent/client/tracker"
	"bittorrent/common"
	"fmt"
	"sync"
)

// **Peer's structure**
type Peer struct {
	Id string

	// Private properties
	address             common.Address
	torrentData         common.Torrent
	notificationChannel chan interface{}
	peers               map[string]PeerInfo // Peers is a <PeerId, PeerInfo> dictionary
	tracker             tracker.Tracker
	fileManager         fileManager.FileManager
	pieceManager        pieceManager.PieceManager
	downloaded          bool
	getAbsoluteOffset   func(int, int) int // function that calculates the absolute offset from index and relative-offset
}

func New(id string, address common.Address, torrent common.Torrent) (Peer, error) {
	peer := Peer{}
	peer.Id = id
	peer.address = address
	peer.torrentData = torrent
	peer.notificationChannel = make(chan interface{}, 1000)
	peer.peers = make(map[string]PeerInfo)

	peer.tracker = tracker.CentralizedHttpTracker{Url: torrent.Announce}

	length := 0
	var files []common.FileInfo
	if torrent.Files == nil {
		files = []common.FileInfo{{
			Length: int(torrent.Length),
			Path:   "./" + torrent.Name,
		}}
		length = int(torrent.Length)
	} else {
		files = torrent.Files
		for _, info := range files {
			length += info.Length
		}
	}

	var err error
	peer.fileManager, err = fileManager.New(files)
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

	go performListen(peer.notificationChannel, peer.address, peer.Id, peer.torrentData.InfoHash)
	go performTrack(peer.notificationChannel, peer.tracker, trackerRequest, 0)
	if !peer.downloaded {
		go performDownload(peer.notificationChannel, 2)
	}

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		for message := range peer.notificationChannel {
			switch notification := message.(type) {
			case killNotification:
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
			default:
				fmt.Println("ERROR: invalid notification-type")
			}
		}
	}()

	waitGroup.Wait()
	return nil
}

func (peer *Peer) handleTrackResponseNotification(notification trackNotification) {
	// TODO: Properly handle the case when the notification was not successful
	fmt.Println("LOG: handling tracker response notification")
	const PEERS_LOWER_BOUND int = 20

	// Check if the file is not downloaded and neighbor-peers are less than the established bound
	if !peer.downloaded && len(peer.peers) < PEERS_LOWER_BOUND {
		for id, address := range notification.Response.Peers {
			if _, contains := peer.peers[id]; !contains {
				go performAddPeer(peer.notificationChannel, peer.Id, id, address, peer.torrentData.InfoHash)
			}
		}
	}

	left := peer.bytesLeft()

	go performTrack(peer.notificationChannel, peer.tracker, common.TrackRequest{
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
	for index := range missing_pieces {
		uncheckedChunks = append(uncheckedChunks, peer.pieceManager.GetUncheckedChunks(index, UNCHECKED_CHUNKS_PER_PIECE)...)
	}

	for _, chunk := range uncheckedChunks {
		for peerId, peerInfo := range peer.peers {
			if peerInfo.Bitfield[chunk[INDEX]] {
				go performSendRequestToPeer(peer.notificationChannel, peerInfo.Connection, peerId, chunk[INDEX], chunk[OFFSET], chunk[LENGTH])
				break
			}
		}
	}

	if !peer.downloaded {
		go performDownload(peer.notificationChannel, 1)
	}
}

func (peer *Peer) handleAddPeerNotification(notification addPeerNotification) {
	_, contains := peer.peers[notification.PeerId]

	if contains {
		return
	}

	peer.peers[notification.PeerId] = PeerInfo{
		Connection: notification.Connection,
		Bitfield:   nil,
		IsChoker:   false,
		IsChoked:   false,
	}
}

func (peer *Peer) handleRemovePeerNotification(notification removePeerNotification) {
	info, contains := peer.peers[notification.PeerId]

	if !contains {
		return
	}

	err := info.Connection.Close()
	if err != nil {
		fmt.Println("ERROR: an error occurred while closing connection " + err.Error())
	}

	delete(peer.peers, notification.PeerId)
	fmt.Println("LOG: remove neighbor: " + notification.PeerId)
}

func (peer *Peer) handleSendBitfieldNotification(notification sendBitfieldNotification) {
	info, contains := peer.peers[notification.PeerId]

	if !contains {
		return
	}

	performSendBitfieldToPeer(peer.notificationChannel, info.Connection, notification.PeerId, peer.pieceManager.Bitfield())
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
	go performSendPieceToPeer(peer.notificationChannel, info.Connection, peer.fileManager, notification.PeerId, notification.Index, notification.Offset, notification.Length, start)
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
	go performWrite(peer.notificationChannel, peer.fileManager, notification.Index, notification.Offset, absoluteOffset, notification.Bytes)

	fmt.Println("LOG: a piece-message was received from: " + notification.PeerId)
}

func (peer *Peer) handleWriteNotification(notification writeNotification) {
	if peer.pieceManager.VerifyPiece(notification.Index) || peer.pieceManager.VerifyChunk(notification.Index, notification.Offset) {
		return
	}

	checkedPiece := peer.pieceManager.CheckChunk(notification.Index, notification.Offset)

	if checkedPiece {
		pieceAbsoluteOffset := notification.Index * int(peer.torrentData.PieceLength)

		go performVerifyPiece(peer.notificationChannel, peer.fileManager, notification.Index, pieceAbsoluteOffset, int(peer.torrentData.PieceLength), peer.torrentData.Pieces)
	}
}

func (peer *Peer) handlePieceVerifiedNotification(notification pieceVerificationNotification) {
	if notification.Verified {
		// TODO: Send a have-message to all peers
		fmt.Printf("LOG: piece %v was verified\n", notification.Index)
	} else {
		fmt.Printf("LOG: piece %v was corrupted\n", notification.Index)
		peer.pieceManager.UncheckPiece(notification.Index)
	}
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
