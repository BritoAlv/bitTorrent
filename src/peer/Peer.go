package peer

import (
	"bittorrent/common"
	"bittorrent/fileManager"
	"bittorrent/tracker"
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

	var files []common.FileInfo
	if torrent.Files == nil {
		files = []common.FileInfo{{
			Length: int(torrent.Length),
			Path:   "./" + torrent.Name,
		}}
	} else {
		files = torrent.Files
	}

	var err error
	peer.fileManager, err = fileManager.New(files)
	if err != nil {
		return peer, err
	}

	peer.getAbsoluteOffset = func(index int, offset int) int {
		return index*int(torrent.PieceLength) + offset
	}
	return peer, nil
}

// **Peer's methods**

func (peer *Peer) Torrent(externalWaitGroup *sync.WaitGroup) error {
	if externalWaitGroup != nil {
		defer externalWaitGroup.Done()
	}

	trackerRequest := tracker.TrackRequest{
		InfoHash: peer.torrentData.InfoHash,
		PeerId:   peer.Id,
		Ip:       peer.address.Ip,
		Port:     peer.address.Port,
		Left:     500,
		// Event:    "started",
	}
	go requestListen(peer.notificationChannel, peer.address)
	go requestTracker(peer.notificationChannel, peer.tracker, trackerRequest, 0)
	go requestDownload(peer.notificationChannel, 10)

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		for message := range peer.notificationChannel {
			switch notification := message.(type) {
			case trackerResponseNotification:
				peer.handleTrackerResponseNotification(notification)
			case downloadNotification:
				peer.handleDownloadNotification()
			case killNotification:
				return
			case peerUpNotification:
				peer.handlePeerUpNotification(notification)
			default:
				fmt.Println("Invalid notification-type")
			}
		}
	}()

	waitGroup.Wait()
	return nil
}

func (peer *Peer) handleTrackerResponseNotification(notification trackerResponseNotification) {
	// TODO: Properly handle the case when the notification was not successful
	fmt.Println("PEER: Handling tracker response notification")
	const PEERS_LOWER_BOUND int = 20

	if len(peer.peers) < PEERS_LOWER_BOUND {
		for id, address := range notification.Response.Peers {
			if _, contains := peer.peers[id]; !contains {
				go requestPeerUp(peer.notificationChannel, id, address)
			}
		}
	}

	go requestTracker(peer.notificationChannel, peer.tracker, tracker.TrackRequest{
		InfoHash: peer.torrentData.InfoHash,
		PeerId:   peer.Id,
		Ip:       peer.address.Ip,
		Port:     peer.address.Port,
		Left:     500,
		// Event:    "started",
	}, notification.Response.Interval)
}

func (peer *Peer) handleDownloadNotification() {
	fmt.Println("PEER: Handling download notification")

	for _, info := range peer.peers {
		message := []byte("Viva Cuba Libre!")
		info.Connection.Write(message)
	}

	go requestDownload(peer.notificationChannel, 10)
}

func (peer *Peer) handlePeerUpNotification(notification peerUpNotification) {
	_, contains := peer.peers[notification.Id]

	if contains {
		return
	}

	peer.peers[notification.Id] = PeerInfo{
		Connection: notification.Connection,
		Bitfield:   nil,
		IsChoker:   false,
		IsChoked:   false,
	}
}
