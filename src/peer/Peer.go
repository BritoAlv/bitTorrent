package peer

import (
	"bittorrent/common"
	"bittorrent/tracker"
	"fmt"
	"sync"
)

// **Peer's structure**
type Peer struct {
	Address             common.Address
	TorrentInfo         TorrentInfo
	Tracker             tracker.Tracker
	NotificationChannel chan interface{}
	Peers               map[common.Address]PeerInfo
}

// **Peer's methods**

func (peer *Peer) Torrent(peerGroup *sync.WaitGroup) error {
	if peerGroup != nil {
		defer peerGroup.Done()
	}

	trackerRequest := tracker.TrackRequest{}
	go requestPeerListen(peer.NotificationChannel, peer.Address)
	go requestTracker(peer.NotificationChannel, peer.Tracker, trackerRequest, 0)
	go requestDownload(peer.NotificationChannel, 10)

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		for message := range peer.NotificationChannel {
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
	fmt.Println("PEER: Handling tracker response notification")
	const PEERS_LOWER_BOUND int = 20

	if len(peer.Peers) < PEERS_LOWER_BOUND {
		for peerAddress, peerId := range notification.Response.Peers {
			if _, contains := peer.Peers[peerAddress]; !contains {
				go requestPeerUp(peer.NotificationChannel, peerAddress, peerId)
			}
		}
	}

	go requestTracker(peer.NotificationChannel, peer.Tracker, tracker.TrackRequest{
		InfoHash: []byte{},
		PeerId:   "",
		Ip:       "",
		Port:     "",
		Left:     0,
		Event:    "",
	}, 5)
}

func (peer *Peer) handleDownloadNotification() {
	fmt.Println("PEER: Handling download notification")

	for _, info := range peer.Peers {
		message := []byte("Viva Cuba Libre!")
		info.Connection.Write(message)
	}

	go requestDownload(peer.NotificationChannel, 10)
}

func (peer *Peer) handlePeerUpNotification(notification peerUpNotification) {
	_, contains := peer.Peers[notification.Address]

	if contains {
		return
	}

	peer.Peers[notification.Address] = PeerInfo{
		Id:         notification.Id,
		Connection: notification.Connection,
		Bitfield:   []bool{},
		IsChoker:   false,
		IsChoked:   false,
	}
}
