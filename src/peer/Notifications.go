package peer

import (
	"bittorrent/common"
	"bittorrent/tracker"
	"net"
)

type trackerResponseNotification struct {
	Response   tracker.TrackResponse
	Successful bool
}

type downloadNotification struct{}

type peerDownNotification struct {
	Address common.Address
}

type peerUpNotification struct {
	Address    common.Address
	Id         string
	Connection net.Conn
}

type killNotification struct{}
