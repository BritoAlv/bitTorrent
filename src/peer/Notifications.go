package peer

import (
	"bittorrent/tracker"
	"net"
)

type trackerResponseNotification struct {
	Response   tracker.TrackResponse
	Successful bool
}

type downloadNotification struct{}

type peerDownNotification struct {
	Id string // Peer's Id
}

type peerUpNotification struct {
	Id         string // Peer's Id
	Connection net.Conn
}

type killNotification struct{}
