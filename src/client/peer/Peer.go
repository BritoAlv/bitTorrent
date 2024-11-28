package peer

import (
	"bittorrent/client/fileManager"
	"bittorrent/client/pieceManager"
	"bittorrent/client/tracker"
	"bittorrent/common"
	"crypto/sha1"
	"fmt"
	"net"
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
	getAbsoluteOffset   func(int, int) int
}

func New(id string, address common.Address, torrent common.Torrent) (Peer, error) {
	peer := Peer{}
	peer.Id = id
	peer.address = address
	peer.torrentData = torrent
	peer.notificationChannel = make(chan interface{}, 1000)
	peer.peers = make(map[string]PeerInfo)

	peer.tracker = tracker.CentralizedTracker{Url: torrent.Announce}

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
		Left:     500,
		// Event:    "started",
	}

	go performListen(peer.notificationChannel, peer.address, peer.Id, peer.torrentData.InfoHash)
	go performTrack(peer.notificationChannel, peer.tracker, trackerRequest, 0)
	if !peer.downloaded {
		go performDownload(peer.notificationChannel, 10)
	}

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		for message := range peer.notificationChannel {
			switch notification := message.(type) {
			case trackNotification:
				peer.handleTrackResponseNotification(notification)
			case downloadNotification:
				peer.handleDownloadNotification()
			case killNotification:
				return
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

	// Calculate the amount of bytes left to download
	left := 0
	if !peer.downloaded {
		for range peer.pieceManager.GetUncheckedPieces() {
			left += int(peer.torrentData.PieceLength)
		}
	}

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
	fmt.Println("LOG: handling download notification")

	// TODO: Properly handle error here
	index, offset, length, _ := peer.pieceManager.GetUncheckedChunk(0)

	var peerId string
	var connection net.Conn
	for id, info := range peer.peers {
		peerId = id
		connection = info.Connection
		break
	}

	if connection != nil {
		go performSendRequestToPeer(peer.notificationChannel, connection, peerId, index, offset, length)
	}

	if !peer.downloaded {
		go performDownload(peer.notificationChannel, 10)
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
}

func (peer *Peer) handleSendBitfieldNotification(notification sendBitfieldNotification) {
	info, contains := peer.peers[notification.PeerId]

	if !contains {
		return
	}

	performSendBitfieldToPeer(peer.notificationChannel, info.Connection, notification.PeerId, peer.pieceManager.Bitfield())
}

func (peer *Peer) handlePeerRequestNotification(notification peerRequestNotification) {
	fmt.Println("LOG: a request-message was received from: " + notification.PeerId)
}

func (peer *Peer) handlePeerBitfieldNotification(notification peerBitfieldNotification) {
	info, contains := peer.peers[notification.PeerId]

	if !contains {
		return
	}

	info.Bitfield = notification.Bitfield
	fmt.Printf("LOG: a bitfield-message was received from: %v -> %v \n", notification.PeerId, notification.Bitfield)
}

func (peer *Peer) checkAllPieces() error {
	for index := range len(peer.torrentData.Pieces) / 20 {
		err := peer.checkPieceHash(index)
		if err != nil {
			return err
		}
	}

	if len(peer.pieceManager.GetUncheckedPieces()) == 0 {
		peer.downloaded = true
	}

	return nil
}

func (peer *Peer) checkPieceHash(index int) error {
	start := peer.getAbsoluteOffset(index, 0)
	bytes, err := peer.fileManager.Read(start, int(peer.torrentData.PieceLength))

	if err != nil {
		// Check if the reading attempt was outside of the file bounds, if so the expected bytes are yet downloaded
		_, isOutsideOfFileBounds := err.(fileManager.OutsideOfFileBoundsError)

		if isOutsideOfFileBounds {
			bytes = []byte{}
		} else {
			return err
		}
	}

	hashIndex := index * 20
	bytesHash := sha1.Sum(bytes)
	pieceHash := peer.torrentData.Pieces[hashIndex : hashIndex+20]

	if bytesHash != [20]byte(pieceHash) {
		peer.pieceManager.UncheckPiece(index)
	}

	return nil
}
